package namespacehelper

import (
	"github.com/devngho/openN-Go/databasehelper"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/types"
	"os"
)

var Namespaces []types.Namespace

func ReadNamespaces() {
	data, err := databasehelper.Dao.ReadNamespaces()
	iohelper.ErrFatal(err)
	Namespaces = data
}

func Find(NamespaceName string) (types.Namespace, error) {
	for _, e := range Namespaces {
		if e.Name == NamespaceName {
			return e, nil
		}
	}
	return types.Namespace{}, os.ErrNotExist
}
