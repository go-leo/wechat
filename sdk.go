package wechat

import (
	"net/http"

	"github.com/go-leo/netx/httpx"
	"github.com/go-leo/stringx"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/go-leo/wechat/auth"
	"github.com/go-leo/wechat/common"
	"github.com/go-leo/wechat/submsg"
)

type options struct {
	httpCli              *http.Client
	appid                string
	secret               string
	redisCli             redis.UniversalClient
	redisSync            *redsync.Redsync
	accessTokenKey       string
	accessTokenLockerKey string
	logger               common.Logger
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *options) init() {
	if o.httpCli == nil {
		o.httpCli = httpx.PooledClient()
	}
	if o.redisCli != nil {
		o.redisSync = redsync.New(goredis.NewPool(o.redisCli))
	}
	if stringx.IsBlank(o.accessTokenKey) {
		o.accessTokenKey = "access:token:" + o.appid
	}
	if stringx.IsBlank(o.accessTokenLockerKey) {
		o.accessTokenLockerKey = o.accessTokenKey + ":locker"
	}
	if o.logger == nil {
		o.logger = &common.DefaultLogger{}
	}
}

type Option func(o *options)

func Logger(l common.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

func HttpClient(httpCli *http.Client) Option {
	return func(o *options) {
		o.httpCli = httpCli
	}
}

func AppID(appid string) Option {
	return func(o *options) {
		o.appid = appid
	}
}

func Secret(secret string) Option {
	return func(o *options) {
		o.secret = secret
	}
}

func RedisClient(client redis.UniversalClient) Option {
	return func(o *options) {
		o.redisCli = client
	}
}

func AccessTokenKey(key string) Option {
	return func(o *options) {
		o.accessTokenKey = key
	}
}

func AccessTokenLockerKey(key string) Option {
	return func(o *options) {
		o.accessTokenLockerKey = key
	}
}

// SDK 微信小程序SDK
type SDK struct {
	o *options
}

func NewSDK(opts ...Option) *SDK {
	o := new(options)
	o.apply(opts...)
	o.init()
	return &SDK{o: o}
}

func (sdk *SDK) AppID() string {
	return sdk.o.appid
}

func (sdk *SDK) Auth() *auth.SDK {
	return &auth.SDK{
		HttpCli:              sdk.o.httpCli,
		AppID:                sdk.o.appid,
		Secret:               sdk.o.secret,
		RedisCli:             sdk.o.redisCli,
		RedisSync:            sdk.o.redisSync,
		AccessTokenKey:       sdk.o.accessTokenKey,
		AccessTokenLockerKey: sdk.o.accessTokenLockerKey,
		Logger:               sdk.o.logger,
	}
}

func (sdk *SDK) SubMsg() *submsg.SDK {
	return &submsg.SDK{
		HttpCli: sdk.o.httpCli,
		AppID:   sdk.o.appid,
	}
}
