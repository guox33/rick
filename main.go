package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func main() {
	str := "{\n    \"Psm\": \"explorer.explorer.broker\",\n    \"FuncName\": \"QueryClusters\",\n    \"ReqBody\": '{\"Psm\": \"toutiao.explorer.explorer\"}',\n    \"IdlSource\": 1,\n    \"IdlVersion\": \"master\",\n    \"Zone\": \"BOEI18N\",\n    \"Idc\": \"boei18n\",\n    \"Cluster\": \"default\",\n    \"Env\": \"prod\",\n    \"Address\": \"\",\n    \"Mock\": \"\",\n    \"LogId\": \"\",\n    \"RequestTimeout\": 1000,\n    \"ConnectTimeout\": 0,\n    \"RequestId\": 0,\n    \"FeatureId\": 0,\n    \"source\": 0,\n    \"Operator\": \"guoxin.rick33\",\n    \"Base\": {\n        \"LogID\": \"\",\n        \"Caller\": \"\",\n        \"Addr\": \"\",\n        \"Client\": \"\",\n        \"TrafficEnv\": {\n            \"Open\": false,\n            \"Env\": \"\"\n        },\n        \"Extra\": {\n            \"\": \"\"\n        }\n    }\n}"
	node := gjson.Get(str, "Base")
	fmt.Println(node.Exists())

	m := map[string]string{
		"user":       "guoxin.rick33",
		"gdpr-token": "a.b.c",
	}

	result, err := sjson.Set(str, "Base.Extra", m)
	fmt.Println(result)
	fmt.Println(err)
}
