package cigbin

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-leo/wechat/common"
)

var (
	URLToken = "https://api.weixin.qq.com/cgi-bin/token"
)

type SDK struct {
	HttpCli        *http.Client
	AppID          string
	Secret         string
	RedisCli       redis.UniversalClient
	RedisSync      *redsync.Redsync
	TokenKey       string
	TokenLockerKey string
	Logger         common.Logger
}
