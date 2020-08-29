package com

import (
	"zlutils/code"
	"zlutils/misc"
)

var (
	ClientErrTodo = code.Add(4101, "todo") //改成你的
	//msg可以是多语言，根据请求头的Device-Language Accept-Language决定语言
	ClientErrTodoMulti = code.Add(4102, code.MSS{
		misc.LangEnglish: "todo",
		misc.LangHindi:   "todo的印地语翻译",
	})
)
