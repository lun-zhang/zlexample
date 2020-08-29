package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
	"zlexample/com"
	"zlexample/model"
	"zlutils/guard"
)

/*
添加书本
@api_doc_tags: app
@api_doc_http_method: POST
@api_doc_relative_paths: /zlexample/app/book/add/:name
*/
func AppAddBook(ctx context.Context, req struct {
	//如果你想要解析uri参数，则成员名必须是Uri，且必须是struct
	Uri struct {
		Name string `uri:"name" json:"name" binding:"required"` //书名
	}
	//如果你想要解析query参数，则成员名必须是Query，且必须是struct
	Query struct {
		Id []int32 `form:"id" json:"id" binding:"required"` //书id
	}
	//如果你想要解析header参数，则成员名必须是Header，且必须是struct
	Header struct {
		GroupId int32 `header:"group_id" json:"group_id"` //条件组，解析错误则为0
	} `bind:"ignore_error"` //正常情况下如果解析发生错误，则会停止解析，并返回err
	//但是有些时候例如条件组参数的时候，发生错误则用默认值时候，就用ignore_error标签

	//如果你想解析body参数，则成员名必须是Body，类型随意，只要能被被json.Unmarshal解析即可
	Body struct {
		Writer model.Writer `json:"writer" binding:"required"` //作者
	}
	//Meta session.Meta
	C *gin.Context `json:"-"` //禁止json编码，因为内部包含了无法编码的成员所以无法打印json日志
}) (resp struct {
	Id      []int32 `json:"id"`       //书id
	Name    string  `json:"name"`     //书名
	GroupId int32   `json:"group_id"` //条件组
	//User    session.User `json:"user"`     //用户
	Ip     string       `json:"ip"`     //ip
	Writer model.Writer `json:"writer"` //书作者
}, err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	entry := logrus.WithContext(ctx).WithField("req", req)

	switch req.Uri.Name {
	case "todo":
		err = com.ClientErrTodo.WithErrorf("name:%s", req.Uri.Name)
		entry.WithError(err).Warn()
		return
	case "multi":
		err = com.ClientErrTodoMulti.WithErrorf("name:%s", req.Uri.Name)
		entry.WithError(err).Warn()
		return
	}
	key := com.ProjectName + ":" + req.Uri.Name
	model.ReConn.GetJson(ctx, key, &resp)
	if len(resp.Id) > 0 {
		return
	}
	//defer model.ReConn.SetJson(ctx, key, &resp, time.Minute*5)//注意，如果defer这样用，就得传&resp
	defer func() {
		model.ReConn.SetJson(ctx, key, resp, time.Minute*5)
	}()

	resp.GroupId = req.Header.GroupId
	resp.Ip = req.C.ClientIP()
	//resp.User = req.Meta.GetUser()
	resp.Id = req.Query.Id
	resp.Name = req.Uri.Name
	resp.Writer = req.Body.Writer
	return
}
