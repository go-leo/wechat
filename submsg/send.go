package submsg

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type SendResp struct {
	common.BaseResp
}

// DataItem 模版内某个 .DATA 的值
type DataItem struct {
	Value any    `json:"value"`
	Color string `json:"color"`
}

// Msg 订阅消息请求参数
type Msg struct {
	ToUser           string               `json:"touser"`            // 必选，接收者（用户）的 openid
	TemplateID       string               `json:"template_id"`       // 必选，所需下发的订阅模板id
	Page             string               `json:"page"`              // 可选，点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Data             map[string]*DataItem `json:"data"`              // 必选, 模板内容
	MiniProgramState string               `json:"miniprogram_state"` // 可选，跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string               `json:"lang"`              // 入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

// Send 发送订阅消息
func (sm *SDK) Send(ctx context.Context, accessToken string, msg *Msg) (*SendResp, error) {
	var resp SendResp
	err := httpx.NewRequestBuilder().
		Post().
		URLString(URLSend).
		Query("access_token", accessToken).
		JSONBody(msg).
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
