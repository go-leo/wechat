package wechat

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/go-leo/wechat/submsg"
)

const appid = ""

const secret = ""

const redisaddr = ""
const redispwd = ""

func TestSDKSubMsgSend(t *testing.T) {
	sdk := NewSDK(
		AppID(appid),
		Secret(secret),
		RedisClient(redis.NewClient(&redis.Options{
			Addr:     redisaddr,
			Password: redispwd,
		})))
	token, err := sdk.Auth().GetAccessToken(context.Background())
	assert.NoError(t, err)
	t.Log(token)
	sendResp, err := sdk.SubMsg().Send(context.Background(), token.AccessToken, &submsg.Msg{
		ToUser:     "oSRk75c5l4Vr1pWBZrg9hPTwDZBs",
		TemplateID: "I6rBMeQD18VrEI4pt8Ng5YLKOPqahr9JjpkwWtRs8jM",
		Page:       fmt.Sprintf("/pages/player/index?id=%d&vId=%d", 1, 1),
		Data: map[string]*submsg.DataItem{
			"thing2": {Value: "标题"},
			"date5":  {Value: "2006-01-02 15:04:05"},
			"thing9": {Value: "正在热播，请及时查看～"},
		},
		MiniProgramState: "trial",
		Lang:             "zh_CN",
	})
	assert.NoError(t, err)
	t.Log(sendResp)
}

func TestSDKSubMsgGetTemplateList(t *testing.T) {
	sdk := NewSDK(
		AppID(appid),
		Secret(secret),
		RedisClient(redis.NewClient(&redis.Options{
			Addr:     redisaddr,
			Password: redispwd,
		})))
	token, err := sdk.Auth().GetAccessToken(context.Background())
	assert.NoError(t, err)
	t.Log(token)
	getTemplateListResp, err := sdk.SubMsg().GetTemplateList(context.Background(), token.AccessToken)
	assert.NoError(t, err)
	t.Log(getTemplateListResp)
}

func TestSDKSubMsgGetPubTemplateTitleList(t *testing.T) {
	sdk := NewSDK(
		AppID(appid),
		Secret(secret),
		RedisClient(redis.NewClient(&redis.Options{
			Addr:     redisaddr,
			Password: redispwd,
		})))
	token, err := sdk.Auth().GetAccessToken(context.Background())
	assert.NoError(t, err)
	t.Log(token)
	getTemplateListResp, err := sdk.SubMsg().GetPubTemplateTitleList(context.Background(), token.AccessToken, []string{""}, 0, 10)

	assert.NoError(t, err)
	t.Log(getTemplateListResp)
}
