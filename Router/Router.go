package Router

import (
	"../DocumentHelper"
	"github.com/gin-gonic/gin"
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
			c.String(404, "Not found")
		}else {
			c.String(200, doc.Name)
		}

		//Document Render
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