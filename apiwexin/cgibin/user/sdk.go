package user

import (
	"net/http"
)

var (
	URLGet = "https://api.weixin.qq.com/cgi-bin/user/get"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
}
