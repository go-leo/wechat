package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-leo/backoffx"
	"github.com/go-leo/mapx"
	"github.com/go-leo/netx/httpx"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-leo/wechat/common"
)

type SDK struct {
	HttpCli              *http.Client
	AppID                string
	Secret               string
	RedisCli             redis.UniversalClient
	RedisSync            *redsync.Redsync
	AccessTokenKey       string
	AccessTokenLockerKey string
	Logger               common.Logger
}

type GetAccessTokenResp struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AppID       string
	AccessToken string `json:"access_token"` // string 获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   // number	凭证有效时间，单位：秒。目前是7200秒之内的值
}

func (auth *SDK) GetAccessToken(ctx context.Context) (*GetAccessTokenResp, error) {
	// 从redis获取token信息，如果获取到了，判断expires_at是否在当前时间之后，如果是之后，则直接返回响应
	// 其他情况需要获取分布式锁，获取到锁的话，就请求微信服务，获取新的token，并更新token，expires_at 设置成的 now+expires_in*(2/3)
	// 没获取到锁的，返回响应（新老有5分钟的过度期）
	result, err := auth.RedisCli.HGetAll(ctx, auth.AccessTokenKey).Result()
	if err != nil {
		return nil, err
	}
	if mapx.IsNotEmpty(result) {
		// 从缓存中获取到token信息
		expiresAtTimestamp, _ := result["expires_at"]
		timestamp, _ := strconv.ParseInt(expiresAtTimestamp, 10, 64)
		expiresAt := time.Unix(timestamp, 0)
		if expiresAt.After(time.Now()) {
			// 没有过期，直接返回token信息
			return auth.decodeResp(result), nil
		} else {
			// 过期了，获取锁
			mutex := auth.RedisSync.NewMutex(auth.AccessTokenLockerKey)
			err := mutex.Lock()
			if err != nil {
				// 获取锁失败，直接返回token信息
				return auth.decodeResp(result), nil
			}
			defer func(mutex *redsync.Mutex) {
				_, _ = mutex.Unlock()
			}(mutex)

			// 获取锁成功, 调微信的接口
			tokenResp, err := auth.getAccessToken(ctx)
			if err != nil {
				auth.Logger.Errorf("failed to get access token from wechat, %v", err)
				return auth.decodeResp(result), nil
			}

			// 保存到redis
			if err := auth.saveToRedis(ctx, tokenResp); err != nil {
				auth.Logger.Errorf("failed to save access token to redis, %v", err)
				return tokenResp, nil
			}
			return tokenResp, nil
		}
	}
	// 从缓存中没有获取到，第一次请求微信获取token
	// 获取锁
	mutex := auth.RedisSync.NewMutex(
		auth.AccessTokenLockerKey,
		redsync.WithTries(3),
		redsync.WithRetryDelayFunc(func(tries int) time.Duration {
			return backoffx.Linear(50*time.Millisecond)(ctx, uint(tries))
		}),
	)
	if err := mutex.Lock(); err != nil {
		// 获取锁失败，在从redis中获取一次
		result, err := auth.RedisCli.HGetAll(ctx, auth.AccessTokenKey).Result()
		if err != nil {
			return nil, err
		}
		if mapx.IsEmpty(result) {
			return nil, errors.New("failed to get access token")
		}
		return auth.decodeResp(result), nil
	}
	defer func(mutex *redsync.Mutex) {
		_, _ = mutex.Unlock()
	}(mutex)
	// 获取锁成功, 调微信的接口
	tokenResp, err := auth.getAccessToken(ctx)
	if err != nil {
		auth.Logger.Errorf("failed to get access token from wechat, %v", err)
		return auth.decodeResp(result), nil
	}

	// 保存到redis
	if err := auth.saveToRedis(ctx, tokenResp); err != nil {
		auth.Logger.Errorf("failed to save access token to redis, %v", err)
		return tokenResp, nil
	}
	return tokenResp, nil

}

func (auth *SDK) saveToRedis(ctx context.Context, tokenResp *GetAccessTokenResp) error {
	data, _ := json.Marshal(tokenResp)
	expiresIn := tokenResp.ExpiresIn * 2 / 3
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	_, err := auth.RedisCli.HMSet(ctx, auth.AccessTokenKey, "resp", string(data), "expires_at", expiresAt.Unix()).Result()
	if err != nil {
		return err
	}
	_, err = auth.RedisCli.ExpireAt(ctx, auth.AccessTokenKey, expiresAt).Result()
	if err != nil {
		return err
	}
	return nil
}

func (auth *SDK) decodeResp(result map[string]string) *GetAccessTokenResp {
	resp, _ := result["resp"]
	authGetAccessTokenResp := &GetAccessTokenResp{}
	_ = json.Unmarshal([]byte(resp), authGetAccessTokenResp)
	return authGetAccessTokenResp
}

func (auth *SDK) getAccessToken(ctx context.Context) (*GetAccessTokenResp, error) {
	req, err := new(httpx.RequestBuilder).
		Get().
		URLString(URLGetAccessToken).
		Query("appid", auth.AppID).
		Query("secret", auth.Secret).
		Query("grant_type", "client_credential").Build(ctx)
	if err != nil {
		return nil, err
	}
	helper := httpx.NewResponseHelper(auth.HttpCli.Do(req))
	if helper.Err() != nil {
		return nil, helper.Err()
	}
	var resp GetAccessTokenResp
	if err := helper.JSONBody(&resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("auth.GetAccessToken error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	return &resp, nil
}

// Code2SessionResp 登录凭证校验的返回结果
type Code2SessionResp struct {
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	AppID      string
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
}

func (auth *SDK) Code2Session(ctx context.Context, code string) (*Code2SessionResp, error) {
	req, err := new(httpx.RequestBuilder).
		Get().
		URLString(URLCode2Session).
		Query("appid", auth.AppID).
		Query("secret", auth.Secret).
		Query("js_code", code).
		Query("grant_type", "authorization_code").Build(ctx)
	if err != nil {
		return nil, err
	}
	helper := httpx.NewResponseHelper(auth.HttpCli.Do(req))
	if helper.Err() != nil {
		return nil, helper.Err()
	}
	var resp Code2SessionResp
	if err := helper.JSONBody(&resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("auth.Code2Session error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = auth.AppID
	return &resp, nil
}
