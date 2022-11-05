package subscribe

import (
	"net/http"
)

var (
	URLSend = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
}
