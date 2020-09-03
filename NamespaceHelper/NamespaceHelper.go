package NamespaceHelper

import (
	"encoding/json"
	"github.com/devngho/openN-Go/ACLHelper"
	"github.com/devngho/openN-Go/IOHelper"
	"os"
	"path/filepath"
)

var Namespaces []Namespace
type Namespace struct {
	Name         string `json:"name"`
	NamespaceACL ACLHelper.ACL `json:"acl"`
}

func ReadNamespaces()  {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	IOHelper.ErrLog(err)
	f, err := os.Open(filepath.Join(dir, "db", "namespaces.json"))
	if os.IsNotExist(err){
		IOHelper.ErrLog(err)
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&Namespaces)
	IOHelper.ErrLog(err)
	_ = f.Close()
}