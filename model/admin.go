package model

import (
	"context"
	"database/sql/driver"
	"github.com/sirupsen/logrus"
	"zlutils/guard"
	"zlutils/mysql"
	"zlutils/session"
)

type Book struct {
	Id                       int32  `gorm:"column:id" json:"id"`         //书id
	Name                     string `gorm:"column:name" json:"name"`     //书名
	Writer                   Writer `gorm:"column:writer" json:"writer"` //作者
	session.OperatorWithTime        //操作者、操作时间
}

func (Book) TableName() string {
	return "book"
}

type Writer struct {
	Name  string `json:"name"`  //姓名
	Phone string `json:"phone"` //手机号
}

func (j Writer) Value() (driver.Value, error) {
	return mysql.Value(j)
}

func (j *Writer) Scan(src interface{}) error {
	return mysql.Scan(j, src)
}

/*
获取书本
@api_doc_tags: admin
@api_doc_http_method: GET
@api_doc_relative_paths: /zlexample/admin/book/get/:id
如果没查到就不返回data字段
*/
func BookGet(ctx context.Context, req struct {
	Uri struct {
		Id int32 `uri:"id" json:"id"` //书本id
	}
}) (resp *Book, err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	entry := logrus.WithContext(ctx).WithField("req", req)

	var bs []Book
	if err = dbConn.
		Where("id = ?", req.Uri.Id).
		WithContext(ctx).
		Find(&bs).
		Error; err != nil {
		entry.WithError(err).Error()
		return
	}
	if len(bs) > 0 {
		return &bs[0], nil
	}
	return
}

/*
获取书本列表
@api_doc_tags: admin
@api_doc_http_method: GET
@api_doc_relative_paths: /zlexample/admin/book/list
*/
func BookList(ctx context.Context, req struct {
	Query struct {
		Page int32 `form:"page" json:"page"` //页号，默认为1
		Size uint8 `form:"size" json:"size"` //页大小，默认为10（最大255）
	}
}) (resp struct {
	Items []Book `json:"items"`
	Total int32  `json:"total"`
}, err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	entry := logrus.WithContext(ctx).WithField("req", req)

	page := int(req.Query.Page)
	size := int(req.Query.Size)

	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	if err = dbConn.
		Table(Book{}.TableName()).
		WithContext(ctx).
		Count(&resp.Total).
		Error; err != nil {
		entry.WithError(err).Error()
		return
	}

	if err = dbConn.
		Order("updated_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		WithContext(ctx).
		Find(&resp.Items).
		Error; err != nil {
		entry.WithError(err).Error()
		return
	}
	return
}

/*
添加书本
@api_doc_tags: admin
@api_doc_http_method: POST
@api_doc_relative_paths: /zlexample/admin/book/add
*/
func BookAdd(ctx context.Context, req struct {
	Body Book
	//Meta session.Meta
}) (err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	entry := logrus.WithContext(ctx).WithField("req", req)

	req.Body.Id = 0 //主键自增
	//req.Body.Operator = req.Meta.GetOperator() //提取操作者
	if err = dbConn.
		WithContext(ctx).
		Create(&req.Body).
		Error; err != nil {
		entry.WithError(err).Error()
		return
	}
	return
}

/*
编辑书本
@api_doc_tags: admin
@api_doc_http_method: POST
@api_doc_relative_paths: /zlexample/admin/book/edit/:id
*/
func BookEdit(ctx context.Context, req struct {
	Uri struct {
		Id int32 `uri:"id" json:"id"` //书本id
	}
	Body Book
	Meta session.Meta
}) (err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	entry := logrus.WithContext(ctx).WithField("req", req)

	req.Body.Id = req.Uri.Id                   //主键自增
	req.Body.Operator = req.Meta.GetOperator() //提取操作者
	if err = dbConn.
		WithContext(ctx).
		Save(&req.Body). //FIXME: 不存在则会创建，不是很符合edit逻辑，不过后台就简单处理了
		Error; err != nil {
		entry.WithError(err).Error()
		return
	}
	return
}

/*
删除书本
@api_doc_tags: admin
@api_doc_http_method: POST
@api_doc_relative_paths: /zlexample/admin/book/delete/:id
*/
func BookDelete(ctx context.Context, req struct {
	Uri struct {
		Id int32 `uri:"id" json:"id"` //书本id
	}
}) (err error) {
	defer guard.BeforeCtx(&ctx)(&err)
	entry := logrus.WithContext(ctx).WithField("req", req)

	//TODO: 可能有的场景需要软删除
	if err = dbConn.
		Where("id = ?", req.Uri.Id).
		WithContext(ctx).
		Delete(&Book{}).
		Error; err != nil {
		entry.WithError(err).Error()
		return
	}
	return
}
