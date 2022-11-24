package component

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type RefreshTokenResp struct {
	common.BaseResp
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

func (sdk *SDK) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResp, error) {
	var resp RefreshTokenResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLAccessToken).
		Query("appid", sdk.AppID).
		Query("grant_type", "refresh_token").
		Query("refresh_token", refreshToken).
		Execute(ctx, sdk.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("subscribeMessage.Send error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sdk.AppID
	return &resp, nil
}
