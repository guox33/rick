package main

import (
	"code.byted.org/gopkg/logs"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
	"io"
	"testing"
	"time"
)

func TestGrpc(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := pb.NewRouteGuideClient(conn)

	rect := &pb.Rectangle{Hi: &pb.Point{Longitude: 1, Latitude: 2}, Lo: &pb.Point{Longitude: 3, Latitude: 4}}
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		panic(err)
	}

	for {
		bytes := proto.MessageV1(make([]byte, 0))
		err = stream.RecvMsg(&bytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(bytes)
	}
}

func SendDeprecatedApiNotifyMessage(c context.Context, user string) error {
	lock.RLock()
	key := getKey("-", "-", user)
	expire, ok := deprecatedNotifyMessageBuffer[key]
	lock.RUnlock()
	if ok && time.Now().Before(expire) {
		return nil
	}
	lock.Lock()
	deprecatedNotifyMessageBuffer[key] = time.Now().Add(1 * time.Hour)
	lock.Unlock()

	elements := []lark2.Element{
		lark2.GenTextElement("我们注意到您正在调用接口测试工具老版本api，接口测试工具v1-v3版本接口未来将不再维护，请尽快将您的服务迁移至v4接口"),
		lark2.GenButtonElement("v4接口使用说明", explorerV4ApiDocUrl),
		lark2.GenButtonElement("加入用户群", larkGroupUrl),
	}

	card := lark2.Card{
		Config: lark2.CardConfig{
			WideScreenMode: false,
		},
		Header: lark2.Header{
			Title: lark2.Title{
				Tag:     "plain_text",
				Content: "接口迁移提醒",
			},
			Template: "yellow",
		},
		Elements: elements,
	}

	req := lark2.SendMessageReq{
		Email:       fmt.Sprintf("%v@bytedance.com", user),
		MsgType:     "interactive",
		UpdateMulti: false,
		Card:        card,
	}

	err := lark2.GetClient().SendMessage(c, req)
	if err != nil {
		return err
	}
	logs.Infof("send message to %d to notify api deprecated", user)
	_ = cmetrics.MetricsClient.EmitCounter("lark.deprecated_notify", 1)
	return nil
}
