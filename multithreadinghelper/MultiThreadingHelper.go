package multithreadinghelper

import (
	"fmt"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/documenthelper"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/markdownhelper"
	"github.com/devngho/openN-Go/namespacehelper"
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

type DocumentCreateRequest struct {
	Name string
	Namespace namespacehelper.Namespace
	Result *[2]string
	StatusCode *int
	UserName string
	WaitChannel *sync.WaitGroup
	Acl string
}

var DocumentReadRequests = make(chan *DocumentReadRequest)
var DocumentCreateRequests = make(chan *DocumentCreateRequest)

func InitGoroutine()  {
	go func() {
		for{
			go ComputeDocumentReadRequest(<-DocumentReadRequests)
		}
	}()
	go func() {
		for{
			go ComputeDocumentCreateRequest(<-DocumentCreateRequests)
		}
	}()
}

func ComputeDocumentReadRequest(req *DocumentReadRequest)  {
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
			doc.Text = markdownhelper.ToHTML(doc.Text)
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
			docHtml = strings.ReplaceAll(docHtml, "${watchacl}", doc.Acl.Watch)
			docHtml = strings.ReplaceAll(docHtml, "${editacl}", doc.Acl.Edit)
			docHtml = strings.ReplaceAll(docHtml, "${editaclacl}", doc.Acl.AclEdit)
			*req.StatusCode = http.StatusForbidden
			*req.Result = [2]string{"text/html; charset=utf-8", docHtml}
		}
	}
	req.WaitChannel.Done()
}

func ComputeDocumentCreateRequest(req *DocumentCreateRequest)  {
	has, err := documenthelper.HasDocument(req.Namespace.Name, req.Name)
	iohelper.ErrLog(err)
	if ! has{
		if aclhelper.AclAllow(req.Acl, req.Namespace.NamespaceACL.Create) {
			_, _ = documenthelper.Create(req.Namespace.Name, req.Name, req.UserName)
			*req.StatusCode = http.StatusFound
			*req.Result = [2]string{"text/html; charset=utf-8", fmt.Sprintf("/w/%s:%s", req.Namespace.Name, req.Name)}
		} else {
			docHtml := themehelper.DocumentAclBlockHtml
			docHtml = strings.ReplaceAll(docHtml, "${namespace}", req.Namespace.Name)
			docHtml = strings.ReplaceAll(docHtml, "${name}", req.Name)
			docHtml = strings.ReplaceAll(docHtml, "${watchacl}", req.Namespace.NamespaceACL.Watch)
			docHtml = strings.ReplaceAll(docHtml, "${editacl}", req.Namespace.NamespaceACL.Edit)
			docHtml = strings.ReplaceAll(docHtml, "${editaclacl}", req.Namespace.NamespaceACL.AclEdit)
			*req.StatusCode = http.StatusForbidden
			*req.Result = [2]string{"text/html; charset=utf-8", docHtml}
		}
	}else{
		docHtml := themehelper.ErrorHtml
		docHtml = strings.ReplaceAll(docHtml, "${error}", "DOCUMENT_ALREADY_EXISTS")
		*req.StatusCode = http.StatusForbidden
		*req.Result = [2]string{"text/html; charset=utf-8", docHtml}
	}
	req.WaitChannel.Done()
}