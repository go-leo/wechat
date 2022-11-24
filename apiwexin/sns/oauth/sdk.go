package oauth

import (
	"net/http"
)

var (
	URLAccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
	Secret  string
}
