package multithreadinghelper

import (
	"fmt"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/documenthelper"
	"github.com/devngho/openN-Go/themehelper"
	"net/http"
	"strings"
	"sync"
)

type DocumentReadRequest struct {
	Name string
	Namespace string
	Result *[2]string
	StatusCode *int
	WaitChannel *sync.WaitGroup
	Acl string
}

var DocumentReadRequests = make(chan *DocumentReadRequest)

func InitGoroutine()  {
	go func() {
		for{
			go ComputeDocumentRequest(<-DocumentReadRequests)
		}
	}()
}

func ComputeDocumentRequest(req *DocumentReadRequest)  {
	doc, err := documenthelper.Read(req.Namespace, req.Name)
	if err != nil{
		docHtml := themehelper.NotFoundDocumentHtml
		docHtml = strings.ReplaceAll(docHtml, "${namespace}", req.Namespace)
		docHtml = strings.ReplaceAll(docHtml, "${name}", req.Name)
		*req.StatusCode = http.StatusNotFound
		*req.Result = [2]string{"text/html; charset=utf-8", docHtml}
	}else {
		//Document Render
		if aclhelper.AclAllow(req.Acl, doc.Acl.Watch) {
			docHtml := themehelper.DocumentHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", doc.Namespace)
			docHtml = strings.ReplaceAll(docHtml, "${name}", doc.Name)
			docHtml = strings.ReplaceAll(docHtml, "${text}", doc.Text)
			docHtml = strings.ReplaceAll(docHtml, "${fname}", fmt.Sprintf("%s:%s", req.Namespace, req.Name))
			*req.StatusCode = http.StatusOK
			*req.Result = [2]string{"text/html; charset=utf-8", docHtml}
		}else{
			docHtml := themehelper.DocumentAclBlockHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", doc.Namespace)
			docHtml = strings.ReplaceAll(docHtml, "${name}", doc.Name)
			*req.StatusCode = http.StatusOK
			*req.Result = [2]string{"text/html; charset=utf-8", docHtml}
		}
	}
	req.WaitChannel.Done()
}