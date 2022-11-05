package newtmpl

import (
	"net/http"
)

var (
	URLGetTemplateList         = "https://api.weixin.qq.com/wxaapi/newtmpl/gettemplate"
	URLGetPubTemplateTitleList = "https://api.weixin.qq.com/wxaapi/newtmpl/getpubtemplatetitles"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
}
