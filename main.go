package main

import (
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/databasehelper"
	"github.com/devngho/openN-Go/markdownhelper"
	"github.com/devngho/openN-Go/namespacehelper"
	"github.com/devngho/openN-Go/router"
	"github.com/devngho/openN-Go/settinghelper"
	"github.com/devngho/openN-Go/themehelper"
	"github.com/devngho/openN-Go/userhelper"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	//Setting
	settinghelper.InitFolderFile()
	settinghelper.LoadSettings()
	databasehelper.SetDAO(settinghelper.ReadSetting("db", "type").String())
	databasehelper.Connection(settinghelper.ReadSetting("db", "server").String(), settinghelper.ReadSetting("db", "setting").String())
	settinghelper.InitData()
	themehelper.InitStatic()
	namespacehelper.ReadNamespaces()
	aclhelper.AclLoad()
	userhelper.Load()
	markdownhelper.SetParser()
	r := gin.Default()

	//Boot server
	r.Use(sessions.Sessions("login", sessions.NewCookieStore([]byte(settinghelper.ReadSetting("secret", "key").String()))))
	router.Setup(r, settinghelper.ReadSetting("wiki", "name").String(), settinghelper.ReadSetting("wiki", "start_page").String())
	r.Run(":80")
}
