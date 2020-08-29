package main

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
	"zlexample/com"
	"zlexample/model"
	"zlutils/caller"
	"zlutils/logger"
	"zlutils/request"
	"zlutils/session"
)

var ctx, _ = xray.BeginSegment(context.Background(), "test")

func init() {
	caller.Init(com.ProjectName)
	logger.Init(logger.Config{Level: logrus.DebugLevel})
}

func TestApp(t *testing.T) {
	req := request.Request{
		Config: request.Config{
			Method: http.MethodPost,
			Url:    "http://localhost:12345/zlexample/app/book/add/math?id=1",
		},
		Query: request.MSI{
			"id": []int32{1, 2, 3},
		},
		Header: request.MSI{
			"group_id":   4,
			"User-Id":    "1",
			"Device-Id":  "d1",
			"Product-Id": session.ProductIdVideoBuddy,
		},
		Body: request.MSI{
			"writer": request.MSI{
				"name":  "名",
				"phone": "123",
			},
		},
	}
	var resp struct {
		request.RespRet
		Data struct {
			Id      []int32      `json:"id"`       //书id
			Name    string       `json:"name"`     //书名
			GroupId int32        `json:"group_id"` //条件组
			User    session.User `json:"user"`     //用户
			Ip      string       `json:"ip"`       //ip
			Writer  model.Writer `json:"writer"`   //书作者
		} `json:"data"`
	}
	if err := req.Do(ctx, &resp); err != nil {
		t.Errorf("app err:%v", err)
	} else {
		t.Logf("app ok, data:%v", resp.Data)
	}
}

func TestAdmin(t *testing.T) {
	req := request.Request{
		Config: request.Config{
			Method: http.MethodPost,
			Url:    "http://localhost:12345/zlexample/admin/book/add",
		},
		Query: request.MSI{
			"user_id":  "1",
			"username": "n1",
		},
		Body: request.MSI{
			"id":   1,
			"name": "math",
			"writer": request.MSI{
				"name":  "名",
				"phone": "123",
			},
		},
	}
	if err := req.Do(ctx, &request.RespRet{}); err != nil {
		t.Errorf("admin err:%v", err)
	} else {
		t.Log("admin ok")
	}
}

func TestOutside(t *testing.T) {
	req := request.Request{
		Config: request.Config{
			Method: http.MethodGet,
			Url:    "http://localhost:12345/zlexample/outside/hi",
		},
	}
	var resp struct {
		request.RespRet
		Data string `json:"data"`
	}
	if err := req.Do(ctx, &resp); err != nil {
		t.Errorf("outside err:%v", err)
	} else {
		t.Logf("outside ok, data:%v", resp.Data)
	}
}

func TestRpc(t *testing.T) {
	req := request.Request{
		Config: request.Config{
			Method: http.MethodGet,
			Url:    "http://localhost:12345/zlexample/rpc/book?caller=test",
		},
	}
	var resp struct {
		request.RespRet
		Data struct {
			Caller string      `json:"caller"` //调用者
			Detail interface{} `json:"detail"`
		} `json:"data"`
	}
	if err := req.Do(ctx, &resp); err != nil {
		t.Errorf("rpc err:%v", err)
	} else {
		t.Logf("rpc ok, data:%v", resp.Data)
	}
}
