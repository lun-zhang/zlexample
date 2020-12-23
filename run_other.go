// +build !windows

package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func run(router *gin.Engine) error {
	logrus.Info("Run Not On windows")
	srv := endless.NewServer(address, router) //endless能优雅重启，但只有*nix上才能用
	return srv.ListenAndServe()
}
