package main

import (
	"./Router"
	"./SettingHelper"
	"github.com/gin-gonic/gin"
)

func main() {
	//Setting
	SettingHelper.InitFolderFile()
	r := gin.Default()

	//PingPong
	r.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong!")
	})

	//Boot server
	Router.Setup(r, SettingHelper.ReadSetting("wiki", "name"))
	r.Run(":80")
}
