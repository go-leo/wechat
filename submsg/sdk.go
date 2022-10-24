package submsg

import (
	"net/http"
)

var (
	URLSend                    = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"
	URLGetTemplateList         = "https://api.weixin.qq.com/wxaapi/newtmpl/gettemplate"
	URLGetPubTemplateTitleList = "https://api.weixin.qq.com/wxaapi/newtmpl/getpubtemplatetitles"
)

type SDK struct {
	HttpCli *http.Client
	AppID   string
}
