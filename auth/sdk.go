package auth

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-leo/wechat/common"
)

var (
	URLCode2Session   = "https://api.weixin.qq.com/sns/jscode2session"
	URLGetAccessToken = "https://api.weixin.qq.com/cgi-bin/token"
	URLGetTicket      = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)

type SDK struct {
	HttpCli              *http.Client
	AppID                string
	Secret               string
	RedisCli             redis.UniversalClient
	RedisSync            *redsync.Redsync
	AccessTokenKey       string
	AccessTokenLockerKey string
	TicketKey            string
	TicketLockerKey      string
	Logger               common.Logger
}
