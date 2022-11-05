package sns

import (
	"net/http"
)

var (
	URLCode2Session = "https://api.weixin.qq.com/sns/jscode2session"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
	Secret  string
}
