package documenthelper

import (
	"context"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/mongohelper"
	"github.com/devngho/openN-Go/namespacehelper"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

type Document struct {
	Namespace string        `json:"namespace"`
	Name      string        `json:"name"`
	Text      string        `json:"text"`
	Acl       aclhelper.ACL `json:"acl"`
	Editor    string        `json:"editor"`
	Version   int           `json:"version"`
}

type DocumentOld struct {
	Namespace string        `json:"namespace"`
	Name      string        `json:"name"`
	Text      string        `json:"text"`
	Acl       aclhelper.ACL `json:"acl"`
	Editor    string        `json:"editor"`
	Version   int           `json:"version"`
	Action    string `json:"action"`
}

const (
	DocumentEditAction = "edit"
	DocumentAclEditAction = "acl_edit"
	DocumentCreateAction = "create"
	DocumentDeleteAction = "delete"
)

func (d Document) Edit(n string){
	
}
func (d Document) Delete() {

}
func Create(namespace string, name string, creater string) (Document, error) {
	_, err := mongohelper.Database.Collection("document").InsertOne(context.TODO(), Document{Namespace: namespace, Name: name, Editor: creater, Version: 1, Text: "", Acl: aclhelper.ACL{
		UseNamespace: true,
	}})
	iohelper.ErrLog(err)
	_, err = mongohelper.Database.Collection("document_old").InsertOne(context.TODO(), bson.M{"namespace":namespace, "name":name, "data":[]DocumentOld{{Namespace: namespace, Name: name, Editor: creater, Version: 1, Action: DocumentCreateAction, Text: "", Acl: aclhelper.ACL{
		UseNamespace: true,
	}}}})
	iohelper.ErrLog(err)
	u, _ := Read(namespace, name)
	return u, err
}
func Read(Namespace string, Name string) (Document, error) {
	var u Document
	res := mongohelper.Database.Collection("document").FindOne(context.TODO(), bson.M{"namespace": Namespace,"name": Name})
	iohelper.ErrLog(res.Err())
	iohelper.ErrLog(res.Decode(&u))
	if u.Acl.UseNamespace{
		n, err := namespacehelper.Find(u.Namespace)
		if err != nil{
			return u, os.ErrInvalid
		}else {
			u.Acl = aclhelper.ACL{Delete: n.NamespaceACL.Delete, Edit: n.NamespaceACL.Edit, AclEdit: n.NamespaceACL.AclEdit, Watch: n.NamespaceACL.Watch, UseNamespace: false}
		}
	}
	return u, nil
}