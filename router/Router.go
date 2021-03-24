package router

import (
	"fmt"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/documenthelper"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/markdownhelper"
	"github.com/devngho/openN-Go/namespacehelper"
	"github.com/devngho/openN-Go/settinghelper"
	"github.com/devngho/openN-Go/themehelper"
	"github.com/devngho/openN-Go/types"
	"github.com/devngho/openN-Go/userhelper"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func OnRequest(c *gin.Context, reqType int8) {
	session := sessions.Default(c)
	//Find types
	switch reqType {
	case types.WatchDocument:
		//Document Name
		temp := strings.Split(c.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], ":")
		if len(temp) == 1 {
			c.Redirect(302, fmt.Sprintf("%s:%s", settinghelper.ReadSetting("default", "namespace"), DocumentNamespace))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\n", DocumentNamespace, DocumentName)
		ns, err := namespacehelper.Find(DocumentNamespace)
		if err != nil {
			DocumentName = fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName)
			DocumentNamespace = settinghelper.ReadSetting("default", "namespace").String()
			ns, _ = namespacehelper.Find(DocumentNamespace)
			c.Redirect(http.StatusFound, fmt.Sprintf("/w/%s:%s", DocumentNamespace, DocumentName))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\nNsName : %s\n", DocumentNamespace, DocumentName, ns.Name)

		//Get Acl
		acl := "ip"
		uid := session.Get("uid")
		if uid != nil {
			usr, err := userhelper.FindUserWithUid(fmt.Sprintf("%v", uid))
			if err != nil {
				session.Clear()
				_ = session.Save()
				c.Redirect(302, "/login")
			} else {
				acl = usr.Acl
			}
		}

		//Document Read
		doc, err := documenthelper.Read(ns.Name, DocumentName)
		if err != nil {
			docHtml := themehelper.NotFoundDocumentHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", ns.Name)
			docHtml = strings.ReplaceAll(docHtml, "${name}", DocumentName)
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(docHtml))
		} else {
			//Document Render
			if aclhelper.AclAllow(acl, doc.Acl.Watch) {
				doc.Text = markdownhelper.ToHTML(doc.Text)
				docHtml := themehelper.DocumentHtml
				docHtml = strings.ReplaceAll(docHtml, "${namespace}", doc.Namespace)
				docHtml = strings.ReplaceAll(docHtml, "${name}", doc.Name)
				docHtml = strings.ReplaceAll(docHtml, "${text}", doc.Text)
				docHtml = strings.ReplaceAll(docHtml, "${fname}", fmt.Sprintf("%s:%s", ns.Name, DocumentName))
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docHtml))
			} else {
				docHtml := themehelper.DocumentAclBlockHtml
				docHtml = strings.ReplaceAll(docHtml, "${namespace}", doc.Namespace)
				docHtml = strings.ReplaceAll(docHtml, "${name}", doc.Name)
				docHtml = strings.ReplaceAll(docHtml, "${watchacl}", doc.Acl.Watch)
				docHtml = strings.ReplaceAll(docHtml, "${editacl}", doc.Acl.Edit)
				docHtml = strings.ReplaceAll(docHtml, "${editaclacl}", doc.Acl.AclEdit)
				c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(docHtml))
			}
		}
	case types.LoginPost:
		id := c.PostForm("id")
		pwd := c.PostForm("password")
		user, err := userhelper.FindUserWithNamePwd(id, sha3.Sum512([]byte(pwd)))
		if err != nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(strings.ReplaceAll(themehelper.LoginHtml, "${error}", "has-error")))
		} else {
			session.Set("uid", user.Uid)
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			} else {
				c.Redirect(302, "/")
			}
		}
	case types.NewDocument:
		temp := strings.Split(c.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], ":")
		if len(temp) == 1 {
			c.Redirect(302, fmt.Sprintf("%s:%s", settinghelper.ReadSetting("default", "namespace"), DocumentNamespace))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\n", DocumentNamespace, DocumentName)

		ns, err := namespacehelper.Find(DocumentNamespace)
		if err != nil {
			DocumentName = fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName)
			DocumentNamespace = settinghelper.ReadSetting("default", "namespace").String()
			ns, _ = namespacehelper.Find(DocumentNamespace)
			c.Redirect(http.StatusFound, fmt.Sprintf("/new/%s:%s", DocumentNamespace, DocumentName))
			return
		}
		fmt.Printf("DocNamespace : %s\nDocName : %s\nNsName : %s\n", DocumentNamespace, DocumentName, ns.Name)

		acl := "ip"
		uid := session.Get("uid")
		var usr types.User
		if uid != nil {
			var err error
			usr, err = userhelper.FindUserWithUid(fmt.Sprintf("%v", uid))
			if err != nil {
				session.Clear()
				_ = session.Save()
				c.Redirect(302, "/login")
			} else {
				acl = usr.Acl
			}
		}

		var creator string
		if acl == "ip" {
			creator = c.ClientIP()
		} else {
			creator = usr.Name
		}

		has, err := documenthelper.HasDocument(ns.Name, DocumentName)
		iohelper.ErrLog(err)
		if !has {
			if aclhelper.AclAllow(acl, ns.NamespaceACL.Create) {
				_, _ = documenthelper.Create(ns.Name, DocumentName, creator)
				c.Redirect(http.StatusFound, fmt.Sprintf("/w/%s:%s", ns.Name, DocumentName))
			} else {
				docHtml := themehelper.DocumentAclBlockHtml
				docHtml = strings.ReplaceAll(docHtml, "${namespace}", ns.Name)
				docHtml = strings.ReplaceAll(docHtml, "${name}", DocumentName)
				docHtml = strings.ReplaceAll(docHtml, "${watchacl}", ns.NamespaceACL.Watch)
				docHtml = strings.ReplaceAll(docHtml, "${editacl}", ns.NamespaceACL.Edit)
				docHtml = strings.ReplaceAll(docHtml, "${editaclacl}", ns.NamespaceACL.AclEdit)
				c.Data(http.StatusForbidden, "text/html; charset=utf-8", []byte(docHtml))
			}
		} else {
			docHtml := themehelper.ErrorHtml
			docHtml = strings.ReplaceAll(docHtml, "${error}", "DOCUMENT_ALREADY_EXISTS")
			c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(docHtml))
		}
	}
}

//Registry Route
func Setup(r *gin.Engine, wikiName string, mainPage string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	r.GET("/w/:document", func(context *gin.Context) {
		OnRequest(context, types.WatchDocument)
	})
	r.GET("/license", func(context *gin.Context) {
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(themehelper.LicenseHtmlFile))
	})
	r.GET("/login", func(context *gin.Context) {
		if sessions.Default(context).Get("uid") == nil {
			context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(strings.ReplaceAll(themehelper.LoginHtml, "${error}", "")))
		} else {
			context.Redirect(302, "/")
		}
	})
	r.POST("/login", func(context *gin.Context) {
		OnRequest(context, types.LoginPost)
	})
	r.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/w/"+wikiName+":"+mainPage)
	})
	r.GET("/acl", func(context *gin.Context) {
		uid := sessions.Default(context).Get("uid")
		if uid != nil {
			t, _ := userhelper.FindUserWithUid(fmt.Sprintf("%v", uid))
			context.JSON(200, gin.H{"Acl": t.Acl, "Includes": aclhelper.AclRoles[t.Acl]})
		}
	})
	r.GET("/new/:document", func(context *gin.Context) {
		temp := strings.Split(context.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], ":")
		if len(temp) == 1 {
			context.Redirect(302, fmt.Sprintf("%s:%s", settinghelper.ReadSetting("default", "namespace"), DocumentNamespace))
			return
		}
		ns, err := namespacehelper.Find(DocumentNamespace)
		if err != nil {
			DocumentName = fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName)
			DocumentNamespace = settinghelper.ReadSetting("default", "namespace").String()
			ns, _ = namespacehelper.Find(DocumentNamespace)
		}
		docHtml := themehelper.DocumentNewHtml
		docHtml = strings.ReplaceAll(docHtml, "${namespace}", ns.Name)
		docHtml = strings.ReplaceAll(docHtml, "${name}", DocumentName)
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docHtml))
	})
	r.GET("/reload", func(context *gin.Context) {
		namespacehelper.ReadNamespaces()
		aclhelper.AclLoad()
		context.String(http.StatusOK, "RELOAD COMPLETED")
	})
	r.POST("/new/:document", func(context *gin.Context) {
		OnRequest(context, types.NewDocument)
	})
	r.Static("/theme/", filepath.Join(dir, "theme"))
	r.StaticFile("/favicon.ico", filepath.Join(dir, "theme", "favicon.ico"))
}
