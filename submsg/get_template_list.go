package submsg

import (
	"context"
	"fmt"
	"net/http"

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
	req, err := new(httpx.RequestBuilder).
		Get().
		URLString(URLGetTemplateList).
		Query("access_token", accessToken).
		Build(ctx)
	if err != nil {
		return nil, err
	}
	helper := httpx.NewResponseHelper(sm.HttpCli.Do(req))
	if helper.Err() != nil {
		return nil, helper.Err()
	}
	if helper.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("subscribeMessage.GetTemplateList, uri=%v , statusCode=%v", req.URL, helper.StatusCode())
	}
	var resp GetTemplateListResp
	if err := helper.JSONBody(&resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("subscribeMessage.GetTemplateList error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sm.AppID
	return &resp, nil
}
