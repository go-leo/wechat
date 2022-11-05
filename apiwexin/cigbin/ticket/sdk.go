package ticket

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-leo/wechat/common"
)

var (
	URLGetTicket = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)

type SDK struct {
	HttpCli            *http.Client
	AppID              string
	Secret             string
	RedisCli           redis.UniversalClient
	RedisSync          *redsync.Redsync
	GetTicketKey       string
	GetTicketLockerKey string
	Logger             common.Logger
}
