package main

import (
	common "code.byted.org/ttarch/byteconf_common"
	openapi "code.byted.org/ttarch/byteconf_openapi"
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"strings"
	"time"
)

const (
	SmartTest = "smart_test"
	ByteQI    = "byteqi"
	User      = "guoxin.rick33"
)

var (
	pathList = []string{
		"domain_setting", "gcp_tiktok_cdaas", "task_schedule",
		"combine_cases", "data_deduplication", "rpc_case_filter",
	}
	pathNameMap = map[string]string{
		"/domain_setting":     "域名配置",
		"/gcp_tiktok_cdaas":   "GCP_TikTok_CDaaS",
		"/task_schedule":      "任务调度服务",
		"/combine_cases":      "组合生成服务",
		"/data_deduplication": "数据去重服务",
		"/rpc_case_filter":    "rpc数据过滤插件",
	}

	schemaNoteMap = map[string]string{
		"domain_setting":            "各个地方域名或者其它参数配置",
		"GCP_TikTok_CDaaS_schema":   "GCP TikTok CDaaS 流水线卡点",
		"task_schedule_schema":      "任务调度服务schema",
		"combine_cases_schema":      "组合生成服务schema",
		"data_deduplication_schema": "数据去重服务schema",
		"rpc_task_filter_schema":    "rpc数据调度服务schema",
	}

	passwordMap = map[openapi.ByteconfRegion]map[string]string{
		openapi.ByteconfRegionCN: {
			SmartTest: "649CDD7E5CA8E58D23AD",
			ByteQI:    "3064026FB7ABBC1863CD",
		},
		openapi.ByteconfRegionUs: {
			SmartTest: "649CDD7E5CA8E58D23AD",
			ByteQI:    "3064026FB7ABBC1863CD",
		},
	}
)

func doMigrateSchema(ctx context.Context, region openapi.ByteconfRegion) (err error) {
	smartTestCli := openapi.NewOpenAPI(User, SmartTest, passwordMap[region][SmartTest],
		openapi.NewOptions().WithRegion(region).WithTimeout(10*time.Second))
	byteqiCli := openapi.NewOpenAPI(User, ByteQI, passwordMap[region][ByteQI],
		openapi.NewOptions().WithRegion(region).WithTimeout(10*time.Second))

	getSchemaResp, err := smartTestCli.WithContext(ctx).GetSchemaList(openapi.GetSchemaListReq{NsKey: SmartTest})

	for _, schema := range getSchemaResp {
		if _, ok := schemaNoteMap[schema.Name]; ok {
			_, err := byteqiCli.WithContext(ctx).GetSchemaInfo(openapi.GetSchemaInfoReq{
				NsKey: ByteQI,
				Name:  schema.Name,
			})
			if err == nil {
				continue
			} else if err.(*openapi.Error).Detail != "获取schema失败, err: record not found" {
				return err
			}

			schemaInstance := common.SchemaNode{}
			schemaBs, _ := json.Marshal(schema.Schema)
			err = json.Unmarshal(schemaBs, &schemaInstance)

			err = byteqiCli.WithContext(ctx).CreateSchema(openapi.CreateSchemaReq{
				NsKey:  "byteqi",
				Name:   schema.Name,
				Note:   schemaNoteMap[schema.Name],
				Schema: &schemaInstance,
			})

			if err != nil {
				return err
			}
		}
	}
	return
}

func doMigrateConfig(ctx context.Context, region openapi.ByteconfRegion) (err error) {
	smartTestCli := openapi.NewOpenAPI(User, SmartTest, passwordMap[region][SmartTest],
		openapi.NewOptions().WithRegion(region).WithTimeout(10*time.Second))
	byteqiCli := openapi.NewOpenAPI(User, ByteQI, passwordMap[region][ByteQI],
		openapi.NewOptions().WithRegion(region).WithTimeout(10*time.Second))

	for _, path := range pathList {
		configList, err := smartTestCli.GetConfigList(openapi.GetConfigListReq{
			NsKey: SmartTest,
			Path:  "/" + path,
		})
		if err != nil {
			return err
		}

		for _, config := range configList {
			if config.BizTreePath == "/task_schedule" && config.Name != "task_schedule_params" {
				continue
			}

			_, err = byteqiCli.GetConfig(openapi.GetConfigReq{
				NsKey: ByteQI,
				Path:  "/" + path,
				Name:  config.Name,
			})
			if err == nil || strings.Contains(err.(*openapi.Error).Detail, "获取配置失败, status_code:11, err:record not found: the config has not been published") {
				continue
			} else if err.(*openapi.Error).Detail != "获取配置失败, status_code:11, err:record not found" {
				return err
			}

			if config.BizTreePath == "/rpc_case_filter" {
				config.SchemaName = "rpc_task_filter_schema"
			}
			getSchema, err := byteqiCli.GetSchemaInfo(openapi.GetSchemaInfoReq{
				NsKey: ByteQI,
				Name:  config.SchemaName,
			})
			if err != nil {
				return err
			}

			content, _ := config.Base.MarshalJSON()
			if config.BizTreePath == "/data_deduplication" {
				tmpMap := make(map[string]interface{})
				_ = json.Unmarshal(content, &tmpMap)
				if _, ok := tmpMap["rpc_body_params_deduplication_filter_list"]; !ok {
					tmpMap["rpc_body_params_deduplication_filter_list"] = []string{}
				}
				content, _ = json.Marshal(tmpMap)
			}
			defaultConf, _ := config.DefaultConf.MarshalJSON()
			var createResp interface{}
			createResp, err = byteqiCli.CreateConfig(openapi.CreateConfigReq{
				Type:         openapi.ConfType(config.Type),
				NsKey:        ByteQI,
				Path:         config.BizTreePath,
				Name:         config.Name,
				DisplayName:  config.DisplayName,
				ConstValList: config.ConstValList,
				Content:      string(content),
				DefaultConf:  string(defaultConf),
				SchemaId:     getSchema.ID,
				Desc:         config.Note,
			})
			if err != nil {
				return err
			}
			fmt.Println(createResp)
		}
	}

	return
}

func doCreateBizPath(ctx context.Context, region openapi.ByteconfRegion) (err error) {
	byteqiCli := openapi.NewOpenAPI(User, ByteQI, passwordMap[region][ByteQI],
		openapi.NewOptions().WithRegion(region).WithTimeout(10*time.Second))

	for path, name := range pathNameMap {
		err = byteqiCli.CreateBizPath(openapi.CreateBizPathReq{
			NsKey: ByteQI,
			Path:  path,
			Name:  name,
			Desc:  "",
		})
		if err != nil {
			return err
		}
	}
	return
}

func doPublishConfig(ctx context.Context, region openapi.ByteconfRegion) {
	byteqiCli := openapi.NewOpenAPI(User, ByteQI, passwordMap[region][ByteQI],
		openapi.NewOptions().WithRegion(region).WithTimeout(60*time.Second))

	for _, path := range pathList {
		getConfigListResp, err := byteqiCli.GetConfigList(openapi.GetConfigListReq{
			NsKey: ByteQI,
			Path:  "/" + path,
		})
		if err != nil {
			panic(err)
		}

		for _, config := range getConfigListResp {
			versionList, err := byteqiCli.GetConfigVersionList(openapi.GetVersionListReq{
				Id:     config.Id,
				Limit:  1,
				Offset: 0,
				NsKey:  ByteQI,
			})
			if err != nil {
				panic(err)
			}
			version := versionList.List[0]

			boeList := []string{"boe"}
			onlineList := []string{"cn", "alisg", "gcp"}
			if region == openapi.ByteconfRegionUs {
				boeList = []string{"boei18n"}
				onlineList = []string{"maliva"}
			}

			err = byteqiCli.PublishConfig(openapi.PublishConfigReq{
				NsKey:                ByteQI,
				Path:                 config.BizTreePath,
				Name:                 config.Name,
				Version:              version.Version,
				PublishOrderType:     openapi.OrderTypeBoe,
				NeedVersionIncrement: false,
				RegionList:           boeList,
			})
			if err != nil {
				fmt.Println(err)
			}

			err = byteqiCli.PublishConfig(openapi.PublishConfigReq{
				NsKey:                ByteQI,
				Path:                 config.BizTreePath,
				Name:                 config.Name,
				Version:              version.Version,
				PublishOrderType:     openapi.OrderTypeOnline,
				PPEChannel:           "",
				NeedVersionIncrement: true,
				RegionList:           onlineList,
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

type Address struct {
	City   string `json:"city"`
	Region string `json:"region"`
}

type People struct {
	Name    string   `json:"name"`
	Age     int32    `json:"age"`
	Address *Address `json:"address"`
}

func (a *Address) GetCity() string {
	return a.City
}

func GetCity(a *Address) string {
	return a.City
}

func main() {
	c := cron.New()
	_ = c.AddFunc("*/12 * * * * *", func() {
		fmt.Println("hello")
	})
	c.Start()
	es := c.Entries()[0]
	fmt.Printf("%+v", es.Next)
	time.Sleep(10 * time.Second)
}
