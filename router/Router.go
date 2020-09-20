package router

import (
	"fmt"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/multithreadinghelper"
	"github.com/devngho/openN-Go/namespacehelper"
	"github.com/devngho/openN-Go/settinghelper"
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
		DocumentName := strings.Join(temp[1:], ":")
		if len(temp) == 1{
			c.Redirect(302, fmt.Sprintf("%s:%s", settinghelper.ReadSetting("default", "namespace"), DocumentNamespace))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\n", DocumentNamespace, DocumentName)
		ns, err := namespacehelper.Find(DocumentNamespace)
		if err != nil{
			DocumentName = fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName)
			DocumentNamespace = settinghelper.ReadSetting("default", "namespace")
			ns, _ = namespacehelper.Find(DocumentNamespace)
			c.Redirect(http.StatusFound, fmt.Sprintf("/w/%s:%s", DocumentNamespace, DocumentName))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\nNsName : %s\n", DocumentNamespace, DocumentName, ns.Name)

		//Get Acl
		acl := "ip"
		uid := session.Get("uid")
		if uid != nil{
			usr, err := userhelper.FindUserWithUid(fmt.Sprintf("%v", uid))
			if err != nil{
				session.Clear()
				_ = session.Save()
				c.Redirect(302, "/login")
			}else{
				acl = usr.Acl
			}
		}

		//Document Read
		var res [2]string
		var statusCode int
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		multithreadinghelper.DocumentReadRequests <- &multithreadinghelper.DocumentReadRequest{Name: DocumentName, Namespace: ns.Name, Result: &res, StatusCode: &statusCode, WaitChannel: &waitGroup, Acl: acl}
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
	case NewDocument:
		temp := strings.Split(c.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], ":")
		if len(temp) == 1{
			c.Redirect(302, fmt.Sprintf("%s:%s", settinghelper.ReadSetting("default", "namespace"), DocumentNamespace))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\n", DocumentNamespace, DocumentName)

		ns, err := namespacehelper.Find(DocumentNamespace)
		if err != nil{
			DocumentName = fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName)
			DocumentNamespace = settinghelper.ReadSetting("default", "namespace")
			ns, _ = namespacehelper.Find(DocumentNamespace)
			c.Redirect(http.StatusFound, fmt.Sprintf("/new/%s:%s", DocumentNamespace, DocumentName))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\nNsName : %s\n", DocumentNamespace, DocumentName, ns.Name)

		acl := "ip"
		uid := session.Get("uid")
		var usr userhelper.User
		if uid != nil{
			var err error
			usr, err = userhelper.FindUserWithUid(fmt.Sprintf("%v", uid))
			if err != nil{
				session.Clear()
				_ = session.Save()
				c.Redirect(302, "/login")
			}else{
				acl = usr.Acl
			}
		}

		var creator string
		if acl == "ip"{
			creator = c.ClientIP()
		}else {
			creator = usr.Name
		}

		var res [2]string
		var statusCode int
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		multithreadinghelper.DocumentCreateRequests <- &multithreadinghelper.DocumentCreateRequest{Name: DocumentName, Namespace: ns.Name, Result: &res, StatusCode: &statusCode, WaitChannel: &waitGroup, Acl: acl, UserName: creator}
		waitGroup.Wait()
		if statusCode == http.StatusFound{
			c.Redirect(302, res[1])
		}else {
			c.Data(statusCode, res[0], []byte(res[1]))
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
			context.JSON(200, gin.H{"Acl": t.Acl, "Includes": aclhelper.AclRoles[t.Acl]})
		}
	})
	r.GET("/new/:document", func(context *gin.Context) {
		temp := strings.Split(context.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], ":")
		if len(temp) == 1{
			context.Redirect(302, fmt.Sprintf("%s:%s", settinghelper.ReadSetting("default", "namespace"), DocumentNamespace))
			return
		}
		ns, err := namespacehelper.Find(DocumentNamespace)
		if err != nil{
			DocumentName = fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName)
			DocumentNamespace = settinghelper.ReadSetting("default", "namespace")
			ns, _ = namespacehelper.Find(DocumentNamespace)
		}
		docHtml := themehelper.DocumentNewHtml
		docHtml = strings.ReplaceAll(docHtml, "${namespace}", ns.Name)
		docHtml = strings.ReplaceAll(docHtml, "${name}", DocumentName)
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docHtml))
	})
	r.POST("/new/:document", func(context *gin.Context) {
		OnRequest(context, NewDocument)
	})
	r.Static("/theme/",filepath.Join(dir, "theme"))
	r.StaticFile("/favicon.ico", filepath.Join(dir, "theme", "favicon.ico"))
}