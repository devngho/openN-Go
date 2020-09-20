package namespacehelper

import (
	"context"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/mongohelper"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var Namespaces []Namespace
type Namespace struct {
	Name         string        `json:"name"`
	NamespaceACL aclhelper.ACLNamespace `json:"acl"`
}

func ReadNamespaces()  {
	var results []Namespace // Use your own type here, but this works too

	cur, err := mongohelper.Database.Collection("namespace").Find(context.TODO(), bson.D{})
	iohelper.ErrLog(err)
	err = cur.All(context.TODO(), &results)
	iohelper.ErrLog(err)
	Namespaces = results
}

func Find(NamespaceName string) (Namespace, error) {
	for _, e := range Namespaces{
		if e.Name == NamespaceName{
			return e, nil
		}
	}
	return Namespace{}, os.ErrNotExist
}