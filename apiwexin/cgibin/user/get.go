package user

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type GetTicketResp struct {
	common.BaseResp
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		Openid []string `json:"openid"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
}

func (auth *SDK) Get(ctx context.Context, accessToken string, NextOpenId string) (*GetTicketResp, error) {
	var resp GetTicketResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLGet).
		Query("access_token", accessToken).
		Query("next_openid", NextOpenId).
		Execute(ctx, auth.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("ticket error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	return &resp, nil
}
