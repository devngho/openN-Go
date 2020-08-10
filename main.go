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

	//Boot server
	Router.Setup(r, SettingHelper.ReadSetting("wiki", "name"), SettingHelper.ReadSetting("wiki", "start_page"))
	r.Run(":80")
}
