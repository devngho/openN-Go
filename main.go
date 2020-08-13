package main

import (
	"github.com/devngho/openN-Go/Router"
	"github.com/devngho/openN-Go/SettingHelper"
	"github.com/devngho/openN-Go/ThemeHelper"
	"github.com/gin-gonic/gin"
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
