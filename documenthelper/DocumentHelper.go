package documenthelper

import (
	"github.com/devngho/openN-Go/databasehelper"
	"github.com/devngho/openN-Go/namespacehelper"
	"github.com/devngho/openN-Go/types"
	"os"
)

func Edit(d types.Document, n string) error {
	return databasehelper.Dao.EditDocument(d, n)
}

func Delete(d types.Document) error {
	return databasehelper.Dao.DeleteDocument(d)
}
func Create(namespace string, name string, creator string) (types.Document, error) {
	err := databasehelper.Dao.CreateDocument(namespace, name, creator)
	if err != nil {
		return types.Document{}, err
	}
	err = databasehelper.Dao.CreateDocumentArchive(namespace, name, creator)
	if err != nil {
		return types.Document{}, err
	}
	u, err := Read(namespace, name)
	return u, err
}
func Read(namespace string, name string) (types.Document, error) {
	u, err := databasehelper.Dao.ReadDocument(namespace, name)
	if err != nil {
		return types.Document{}, err
	}
	if u.Acl.UseNamespace {
		n, err := namespacehelper.Find(u.Namespace)
		if err != nil {
			return u, os.ErrInvalid
		} else {
			u.Acl = types.ACL{Delete: n.NamespaceACL.Delete, Edit: n.NamespaceACL.Edit, AclEdit: n.NamespaceACL.AclEdit, Watch: n.NamespaceACL.Watch, UseNamespace: false}
		}
	}
	return u, nil
}
func HasDocument(namespace string, name string) (bool, error) {
	return databasehelper.Dao.HasDocument(namespace, name)
}
