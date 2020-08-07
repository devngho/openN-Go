package main

import (
	"github.com/gin-gonic/gin"
	"openN-Go/Router"
	"openN-Go/SettingHelper"
	"openN-Go/ThemeHelper"
)

func main() {
	//Setting
	SettingHelper.InitFolderFile()
	ThemeHelper.InitStatic()
	r := gin.Default()

	//PingPong
	r.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong!")
	})

	//Boot server
	Router.Setup(r, SettingHelper.ReadSetting("wiki", "name"), SettingHelper.ReadSetting("wiki", "start_page"))
	r.Run(":80")
}
