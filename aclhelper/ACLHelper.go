package aclhelper

import (
	"github.com/devngho/openN-Go/databasehelper"
	"github.com/devngho/openN-Go/iohelper"
)

var AclRoles = make(map[string][]string)

func AclLoad() {
	tmp, err := databasehelper.Dao.ReadACLRole()
	iohelper.ErrFatal(err)
	for _, e := range tmp {
		AclRoles[e.Name] = e.Include
	}
}

func AclAllow(a string, b string) bool {
	if a == b {
		return true
	}
	for _, e := range AclRoles[a] {
		if e == b {
			return true
		}
	}
	return false
}
