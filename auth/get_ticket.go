package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-leo/backoffx"
	"github.com/go-leo/mapx"
	"github.com/go-leo/netx/httpx"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-leo/wechat/common"
)

type GetTicketResp struct {
	common.BaseResp
	Ticket    string `json:"ticket"`     // string 临时票据，用于在获取授权链接时作为参数传入
	ExpiresIn int    `json:"expires_in"` // number	凭证有效时间，单位：秒。目前是7200秒之内的值
}

func (auth *SDK) GetTicket(ctx context.Context, accessToken string, typ string) (*GetTicketResp, error) {
	key := auth.TicketKey + ":" + typ
	lockerKey := auth.TicketLockerKey + ":" + typ
	result, err := auth.RedisCli.HGetAll(ctx, key).Result()
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
			return auth.decodeGetTicketResp(result), nil
		} else {
			// 过期了，获取锁
			mutex := auth.RedisSync.NewMutex(lockerKey)
			err := mutex.Lock()
			if err != nil {
				// 获取锁失败，直接返回token信息
				return auth.decodeGetTicketResp(result), nil
			}
			defer func(mutex *redsync.Mutex) {
				_, _ = mutex.Unlock()
			}(mutex)

			// 获取锁成功, 调微信的接口
			ticketResp, err := auth.getTicket(ctx, accessToken, typ)
			if err != nil {
				auth.Logger.Errorf("failed to get access token from wechat, %v", err)
				return auth.decodeGetTicketResp(result), nil
			}

			// 保存到redis
			if err := auth.saveGetTicketRespToRedis(ctx, key, ticketResp); err != nil {
				auth.Logger.Errorf("failed to save access token to redis, %v", err)
				return ticketResp, nil
			}
			return ticketResp, nil
		}
	}
	// 从缓存中没有获取到，第一次请求微信获取token
	// 获取锁
	mutex := auth.RedisSync.NewMutex(
		lockerKey,
		redsync.WithTries(3),
		redsync.WithRetryDelayFunc(func(tries int) time.Duration {
			return backoffx.Linear(50*time.Millisecond)(ctx, uint(tries))
		}),
	)
	if err := mutex.Lock(); err != nil {
		// 获取锁失败，在从redis中获取一次
		result, err := auth.RedisCli.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if mapx.IsEmpty(result) {
			return nil, errors.New("failed to get ticket")
		}
		return auth.decodeGetTicketResp(result), nil
	}
	defer func(mutex *redsync.Mutex) {
		_, _ = mutex.Unlock()
	}(mutex)

	// 获取锁成功, 调微信的接口
	tokenResp, err := auth.getTicket(ctx, accessToken, typ)
	if err != nil {
		auth.Logger.Errorf("failed to get ticket from wechat, %v", err)
		return auth.decodeGetTicketResp(result), nil
	}

	// 保存到redis
	if err := auth.saveGetTicketRespToRedis(ctx, key, tokenResp); err != nil {
		auth.Logger.Errorf("failed to save access token to redis, %v", err)
		return tokenResp, nil
	}
	return tokenResp, nil
}

func (auth *SDK) saveGetTicketRespToRedis(ctx context.Context, key string, tokenResp *GetTicketResp) error {
	data, _ := json.Marshal(tokenResp)
	expiresIn := tokenResp.ExpiresIn * 2 / 3
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	_, err := auth.RedisCli.HMSet(ctx, key, "resp", string(data), "expires_at", expiresAt.Unix()).Result()
	if err != nil {
		return err
	}
	_, err = auth.RedisCli.ExpireAt(ctx, key, expiresAt).Result()
	if err != nil {
		return err
	}
	return nil
}

func (auth *SDK) decodeGetTicketResp(result map[string]string) *GetTicketResp {
	resp, _ := result["resp"]
	getTicketResp := &GetTicketResp{}
	_ = json.Unmarshal([]byte(resp), getTicketResp)
	return getTicketResp
}

func (auth *SDK) getTicket(ctx context.Context, accessToken string, typ string) (*GetTicketResp, error) {
	var resp GetTicketResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLGetTicket).
		Query("access_token", accessToken).
		Query("type", typ).
		Execute(ctx, auth.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("auth.GetAccessToken error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	return &resp, nil
}
