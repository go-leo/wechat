package urlscheme

import "net/http"

const (
	URLGenerate = "https://api.weixin.qq.com/wxa/generatescheme"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
}
