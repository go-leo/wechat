package newtmpl

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-leo/netx/httpx"

	"github.com/go-leo/wechat/common"
)

type GetTemplateListData struct {
	Tid        int    `json:"tid"`
	Title      string `json:"title"`
	Type       int    `json:"type"`
	CategoryId string `json:"categoryId"`
}

type GetPubTemplateTitleListResp struct {
	common.BaseResp
	Data []*GetTemplateListData `json:"data"`
}

// GetPubTemplateTitleList 获取帐号所属类目下的公共模板标题
func (sm *SDK) GetPubTemplateTitleList(ctx context.Context, accessToken string, ids []string, start int, limit int) (*GetPubTemplateTitleListResp, error) {
	var resp GetPubTemplateTitleListResp
	err := httpx.NewRequestBuilder().
		Get().
		URLString(URLGetPubTemplateTitleList).
		Query("access_token", accessToken).
		Query("ids", strings.Join(ids, ",")).
		Query("start", strconv.Itoa(start)).
		Query("limit", strconv.Itoa(limit)).
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
