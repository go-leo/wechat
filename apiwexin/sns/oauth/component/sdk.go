package component

import (
	"net/http"
)

var (
	URLAccessToken = "https://api.weixin.qq.com/sns/oauth2/component/access_token"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
}
