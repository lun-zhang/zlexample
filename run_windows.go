// +build windows

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func run(router *gin.Engine) error {
	logrus.Info("Run On windows")
	return router.Run(address)
}
