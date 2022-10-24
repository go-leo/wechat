package submsg

import (
	"context"
	"fmt"
	"net/http"
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
	req, err := new(httpx.RequestBuilder).
		Get().
		URLString(URLGetPubTemplateTitleList).
		Query("access_token", accessToken).
		Query("ids", strings.Join(ids, ",")).
		Query("start", strconv.Itoa(start)).
		Query("limit", strconv.Itoa(limit)).
		Build(ctx)
	if err != nil {
		return nil, err
	}
	helper := httpx.NewResponseHelper(sm.HttpCli.Do(req))
	if helper.Err() != nil {
		return nil, helper.Err()
	}
	if helper.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("subscribeMessage.Send, uri=%v , statusCode=%v", req.URL, helper.StatusCode())
	}
	var resp GetPubTemplateTitleListResp
	if err := helper.JSONBody(&resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		err = fmt.Errorf("subscribeMessage.Send error : errcode=%v , errmsg=%v", resp.ErrCode, resp.ErrMsg)
		return nil, err
	}
	resp.AppID = sm.AppID
	return &resp, nil
}
