package submsg

import (
	"context"
	"fmt"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type KeywordEnumValue struct {
	EnumValueList []string `json:"enumValueList"`
	KeywordCode   string   `json:"keywordCode"`
}

type GetTemplateListRespData struct {
	PriTmplId            string              `json:"priTmplId"`
	Title                string              `json:"title"`
	Content              string              `json:"content"`
	Example              string              `json:"example"`
	Type                 int                 `json:"type"`
	KeywordEnumValueList []*KeywordEnumValue `json:"keywordEnumValueList,omitempty"`
}

type GetTemplateListResp struct {
	common.BaseResp
	Data []*GetTemplateListRespData `json:"data"`
}

// GetTemplateList 获取当前帐号下的个人模板列表
func (sm *SDK) GetTemplateList(ctx context.Context, accessToken string) (*GetTemplateListResp, error) {
	var resp GetTemplateListResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLGetTemplateList).
		Query("access_token", accessToken).
		Execute(ctx, sm.HttpCli).
		JSONBody(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("subscribeMessage.GetTemplateList error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sm.AppID
	return &resp, nil
}
