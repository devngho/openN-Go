package Router

import (
	"../DocumentHelper"
	"../ThemeHelper"
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

func OnRequest(c *gin.Context, reqType int8, wikiName string){
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
			c.String(http.StatusNotFound, "Not found")
		}else {
			//Document Render
			docHtml := ThemeHelper.DocumentHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", doc.Namespace)
			docHtml = strings.ReplaceAll(docHtml, "${name}", doc.Name)
			docHtml = strings.ReplaceAll(docHtml, "${text}", doc.Text)
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docHtml))
		}
	}
}

//Registry Route
func Setup(r *gin.Engine, wikiName string){
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	r.GET("/w/:document", func(context *gin.Context) {
		OnRequest(context, WatchDocument, wikiName)
	})
	r.Static("/theme/",filepath.Join(dir, "theme"))
}