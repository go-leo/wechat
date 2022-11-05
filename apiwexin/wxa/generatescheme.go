package wxa

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type JumpWxa struct {
	// Path 通过 scheme 码进入的小程序页面路径，必须是已经发布的小程序存在的页面，不可携带 query。path 为空时会跳转小程序主页。
	Path string `json:"path"`
	// Query 通过 scheme 码进入小程序时的 query，最大1024个字符，只支持数字，大小写英文以及部分特殊字符：`!#$&'()*+,/:;=?@-._~%``
	Query string `json:"query"`
	// EnvVersion 要打开的小程序版本。正式版为"release"，体验版为"trial"，开发版为"develop"，仅在微信外打开时生效。
	EnvVersion string `json:"env_version"`
}

type GenerateSchemeReq struct {
	// JumpWxa 跳转到的目标小程序信息。
	JumpWxa *JumpWxa `json:"jump_wxa"`
	// 默认值false。生成的 scheme 码类型，到期失效：true，永久有效：false。注意，永久有效 scheme 和有效时间超过180天的到期失效 scheme 的总数上限为10万个，详见获取 URL scheme，生成 scheme 码前请仔细确认。
	IsExpire bool `json:"is_expire"`
	// ExpireType 到期失效的 scheme 码失效类型，失效时间：0，失效间隔天数：1
	ExpireType int32 `json:"expire_type"`
	// ExpireTime 到期失效的 scheme 码的失效时间，为 Unix 时间戳。生成的到期失效 scheme 码在该时间前有效。最长有效期为30天。expire_type 为 0 时必填
	ExpireTime int64 `json:"expire_time"`
	// ExpireInterval 到期失效的 scheme 码的失效间隔天数。生成的到期失效 scheme 码在该间隔时间到达前有效。最长间隔天数为30天。 expire_type 为 1 时必填
	ExpireInterval int32 `json:"expire_interval"`
}

type GenerateSchemeResp struct {
	common.BaseResp
	Openlink string `json:"openlink"`
}

// GenerateScheme 获取小程序 scheme 码，适用于短信、邮件、外部网页、微信内等拉起小程序的业务场景
func (sm *SDK) GenerateScheme(ctx context.Context, accessToken string, req *GenerateSchemeReq) (*GenerateSchemeResp, error) {
	var resp GenerateSchemeResp
	err := httpx.NewRequestBuilder().
		Post().
		URLString(URLGenerate).
		Query("access_token", accessToken).
		JSONBody(req).
		Execute(ctx, sm.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("wxa.GenerateScheme error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sm.AppID
	return &resp, nil
}
