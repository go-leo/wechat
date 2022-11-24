package component

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type AccessTokenReq struct {
	Code                 string
	ComponentAppID       string
	ComponentAccessToken string
}

type AccessTokenResp struct {
	common.BaseResp
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

func (sdk *SDK) AccessToken(ctx context.Context, req *AccessTokenReq) (*AccessTokenResp, error) {
	var resp AccessTokenResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLAccessToken).
		Query("appid", sdk.AppID).
		Query("code", req.Code).
		Query("grant_type", "authorization_code").
		Query("component_appid", req.ComponentAppID).
		Query("component_access_token", req.ComponentAccessToken).
		Execute(ctx, sdk.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("AccessToken error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sdk.AppID
	return &resp, nil
}
