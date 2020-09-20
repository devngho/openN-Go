package main

import (
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/mongohelper"
	"github.com/devngho/openN-Go/multithreadinghelper"
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
	mongohelper.Connect(settinghelper.ReadSetting("db", "server"), settinghelper.ReadSetting("db", "database"))
	settinghelper.InitData()
	themehelper.InitStatic()
	namespacehelper.ReadNamespaces()
	aclhelper.AclLoad()
	multithreadinghelper.InitGoroutine()
	userhelper.Load()
	r := gin.Default()

	//Boot server
	r.Use(sessions.Sessions("login", sessions.NewCookieStore([]byte(settinghelper.ReadSetting("secret", "key")))))
	router.Setup(r, settinghelper.ReadSetting("wiki", "name"), settinghelper.ReadSetting("wiki", "start_page"))
	r.Run(":80")
}
