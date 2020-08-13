package Router

import (
	"fmt"
	"github.com/devngho/openN-Go/DocumentHelper"
	"github.com/devngho/openN-Go/ThemeHelper"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const WatchDocument = 0
const EditDocument = 1
const WatchOldlistDocument = 2
const WatchRevDocument = 3
const Search = 4
const NewDocument = 5
const DeleteDocument = 6
const EditAclDocument = 7

func OnRequest(c *gin.Context, reqType int8) {
	//Find type
	switch reqType {
	case WatchDocument:
		//Document Name
		temp := strings.Split(c.Param("document"), ":")
		DocumentNamespace := temp[0]
		DocumentName := strings.Join(temp[1:], "")

		//Document Read
		doc, err := DocumentHelper.Read(DocumentNamespace, DocumentName)
		if err != nil{
			docHtml := ThemeHelper.NotFoundDocumentHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", DocumentNamespace)
			docHtml = strings.ReplaceAll(docHtml, "${name}", DocumentName)
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docHtml))
		}else {
			//Document Render
			docHtml := ThemeHelper.DocumentHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", doc.Namespace)
			docHtml = strings.ReplaceAll(docHtml, "${name}", doc.Name)
			docHtml = strings.ReplaceAll(docHtml, "${text}", doc.Text)
			docHtml = strings.ReplaceAll(docHtml, "${fname}", fmt.Sprintf("%s:%s", DocumentNamespace, DocumentName))
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docHtml))
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
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(ThemeHelper.LicenseHtmlFile))
	})
	r.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/w/"+wikiName+":"+mainPage)
	})
	r.Static("/theme/",filepath.Join(dir, "theme"))
	r.StaticFile("/favicon.ico", filepath.Join(dir, "theme", "favicon.ico"))
}