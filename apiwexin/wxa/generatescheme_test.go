package wxa

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestName(t *testing.T) {
	//a := `{"jump_wxa":{"path":"pages/index/index","query":"","env_version":""},"is_expire":true,"expire_type":0,"expire_time":1667308083184,"expire_interval":0}`

	data, _ := json.Marshal(&GenerateSchemeReq{
		JumpWxa: &JumpWxa{
			Path:       "pages/index/index",
			Query:      "id=10001&a=3",
			EnvVersion: "",
		},
		IsExpire:       false,
		ExpireType:     0,
		ExpireTime:     0,
		ExpireInterval: 0,
	})

	data2, _ := jsoniter.Marshal(&GenerateSchemeReq{
		JumpWxa: &JumpWxa{
			Path:       "pages/index/index",
			Query:      "id=10001&a=3",
			EnvVersion: "",
		},
		IsExpire:       false,
		ExpireType:     0,
		ExpireTime:     0,
		ExpireInterval: 0,
	})

	body := string(data)

	body = strings.Replace(body, `\u0026`, "&", -1)

	fmt.Println(data, data2)
}
