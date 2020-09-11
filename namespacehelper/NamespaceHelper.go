package namespacehelper

import (
	"encoding/json"
	"github.com/devngho/openN-Go/aclhelper"
	"github.com/devngho/openN-Go/iohelper"
	"os"
	"path/filepath"
)

var Namespaces []Namespace
type Namespace struct {
	Name         string        `json:"name"`
	NamespaceACL aclhelper.ACL `json:"acl"`
}

func ReadNamespaces()  {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	iohelper.ErrLog(err)
	f, err := os.Open(filepath.Join(dir, "db", "namespaces.json"))
	if os.IsNotExist(err){
		iohelper.ErrLog(err)
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&Namespaces)
	iohelper.ErrLog(err)
	_ = f.Close()
}