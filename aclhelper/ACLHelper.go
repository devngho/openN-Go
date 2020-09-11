package aclhelper

import (
	"encoding/json"
	"github.com/devngho/openN-Go/iohelper"
	"os"
	"path/filepath"
)

type ACL struct {
	Watch string `json:"watch"`
	Edit string `json:"edit"`
	AclEdit string `json:"acl_edit"`
}

type ACLRole struct {
	Name string `json:"name"`
	Include []string `json:"include"`
}

var AclRoles = make(map[string][]string)

func AclLoad(){
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	iohelper.ErrLog(err)
	f, err := os.Open(filepath.Join(dir, "db", "acl.json"))
	if os.IsNotExist(err){
		iohelper.ErrLog(err)
	}
	dec := json.NewDecoder(f)
	var tmp []ACLRole
	err = dec.Decode(&tmp)
	iohelper.ErrLog(err)
	_ = f.Close()
	for _, e := range tmp{
		AclRoles[e.Name] = e.Include
	}
}

func AclAllow(a string, b string) bool{
	for _, e := range AclRoles[a]{
		if e == b{
			return true
		}
	}
	return false
}
