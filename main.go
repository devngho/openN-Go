package main

import (
	"./Router"
	"./SettingHelper"
	"./ThemeHelper"
	"github.com/gin-gonic/gin"
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
