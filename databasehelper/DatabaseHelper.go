package databasehelper

import (
	"errors"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/types"
)

type Database interface {
	CreateDocument(namespace string, name string, creator string) error
	CreateDocumentArchive(namespace string, name string, creator string) error
	EditDocument(document types.Document, text string) error
	ReadDocument(namespace string, name string) (types.Document, error)
	HasDocument(namespace string, name string) (bool, error)
	DeleteDocument(document types.Document) error
	AclEditDocument(document types.Document, acl types.ACL) error
	ReadACLRole() ([]types.ACLRole, error)
	ReadNamespaces() ([]types.Namespace, error)
	ReadUsers() ([]types.User, error)
	SaveUsers([]types.User) error
	InitData() error
	Connect(connection string, setting string) error
}

var daoList = make(map[string]Database)
var Dao Database

func register(name string, database Database) {
	daoList[name] = database
}

func SetDAO(name string) {
	if daoList[name] == nil {
		iohelper.ErrFatal(errors.New("can't found database " + name))
	} else {
		Dao = daoList[name]
	}
}

func Connection(connection string, setting string) {
	iohelper.ErrFatal(Dao.Connect(connection, setting))
}
