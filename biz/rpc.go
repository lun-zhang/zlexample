package biz

import (
	"context"
	"net/http"
	"zlutils/guard"
	"zlutils/request"
	"zlutils/session"
)

/*
转调app接口
@api_doc_tags: rpc
@api_doc_http_method: GET
@api_doc_relative_paths: /zlexample/rpc/book
*/
func RpcGetBook(ctx context.Context, req struct {
	Query struct {
		//TODO: rpc调用者放到session包中作为中间件
		Caller string `form:"caller" json:"caller" binding:"required"` //调用者
	}
}) (resp struct {
	Caller string      `json:"caller"` //调用者
	Detail interface{} `json:"detail"`
}, err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	resp.Caller = req.Query.Caller

	rq := request.Request{
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
	var rp struct {
		request.RespRet
		Data interface{} `json:"data"`
	}
	if err = rq.Do(ctx, &rp); err != nil {
		return
	}
	resp.Detail = rp.Data
	return
}
