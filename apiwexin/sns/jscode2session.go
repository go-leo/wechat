package sns

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

// JsCode2SessionResp 登录凭证校验的返回结果
type JsCode2SessionResp struct {
	common.BaseResp
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
}

func (auth *SDK) JsCode2Session(ctx context.Context, code string) (*JsCode2SessionResp, error) {
	var resp JsCode2SessionResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLJsCode2Session).
		Query("appid", auth.AppID).
		Query("secret", auth.Secret).
		Query("js_code", code).
		Query("grant_type", "authorization_code").
		Execute(ctx, auth.HttpCli).JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err := fmt.Errorf("sns.JsCode2Session error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = auth.AppID
	return &resp, nil
}
