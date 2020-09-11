package router

import (
	"fmt"
	"github.com/devngho/openN-Go/multithreadinghelper"
	"github.com/devngho/openN-Go/themehelper"
	"github.com/devngho/openN-Go/userhelper"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const WatchDocument = 0
const EditDocument = 1
const WatchOldlistDocument = 2
const WatchRevDocument = 3
const Search = 4
const NewDocument = 5
const DeleteDocument = 6
const EditAclDocument = 7
const LoginPost = 8

func OnRequest(c *gin.Context, reqType int8) {
	session := sessions.Default(c)
	//Find type
	switch reqType {
	case WatchDocument:
		//Document Name
		temp := strings.Split(c.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], "")

		//Document Read
		var res [2]string
		var statusCode int
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		multithreadinghelper.DocumentReadRequests <- &multithreadinghelper.DocumentReadRequest{Name: DocumentName, Namespace: DocumentNamespace, Result: &res, StatusCode: &statusCode, WaitChannel: &waitGroup}
		waitGroup.Wait()
		c.Data(statusCode, res[0], []byte(res[1]))
	case LoginPost:
		id := c.PostForm("id")
		pwd := c.PostForm("password")
		user, err := userhelper.FindUserWithNamePwd(id, sha3.Sum512([]byte(pwd)))
		if err != nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(strings.ReplaceAll(themehelper.LoginHtml, "${error}", "has-error")))
		}else {
			session.Set("uid", user.Uid)
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			} else {
				c.Redirect(302, "/")
			}
		}
	}
}

//Registry Route
func Setup(r *gin.Engine, wikiName string, mainPage string){
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	r.GET("/w/:document", func(context *gin.Context) {
		OnRequest(context, WatchDocument)
	})
	r.GET("/license", func(context *gin.Context) {
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(themehelper.LicenseHtmlFile))
	})
	r.GET("/login", func(context *gin.Context) {
		if sessions.Default(context).Get("uid") == nil {
			context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(strings.ReplaceAll(themehelper.LoginHtml, "${error}", "")))
		}else{
			context.Redirect(302, "/")
		}
	})
	r.POST("/login", func(context *gin.Context) {
		OnRequest(context, LoginPost)
	})
	r.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/w/"+wikiName+":"+mainPage)
	})
	r.GET("/acl", func(context *gin.Context) {
		uid := sessions.Default(context).Get("uid")
		if uid != nil{
			t, _ := userhelper.FindUserWithUid(fmt.Sprintf("%v", uid))
			context.JSON(200, gin.H{"Acl": t.Acl.Name, "Includes": t.Acl.Include})
		}
	})
	r.Static("/theme/",filepath.Join(dir, "theme"))
	r.StaticFile("/favicon.ico", filepath.Join(dir, "theme", "favicon.ico"))
}