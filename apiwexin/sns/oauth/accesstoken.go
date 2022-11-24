package oauth

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type AccessTokenResp struct {
	common.BaseResp
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	AppID        string
}

func (sm *SDK) AccessToken(ctx context.Context, code string) (*AccessTokenResp, error) {
	var resp AccessTokenResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLAccessToken).
		Query("appid", sm.AppID).
		Query("secret", sm.Secret).
		Query("code", code).
		Query("grant_type", "authorization_code").
		Execute(ctx, sm.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("subscribeMessage.Send error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sm.AppID
	return &resp, nil
}
