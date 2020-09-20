package aclhelper

import (
	"context"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/mongohelper"
	"go.mongodb.org/mongo-driver/bson"
)

type ACL struct {
	Watch string `json:"watch"`
	Edit string `json:"edit"`
	AclEdit string `json:"acl_edit"`
	Delete string `json:"delete"`
	UseNamespace bool `json:"use_namespace"`
}

type ACLNamespace struct {
	Watch string `json:"watch"`
	Edit string `json:"edit"`
	AclEdit string `json:"acl_edit"`
	Create string `json:"create"`
	Delete string `json:"delete"`
	UseNamespace bool `json:"use_namespace"`
}

type ACLRole struct {
	Name string `json:"name"`
	Include []string `json:"include"`
}

var AclRoles = make(map[string][]string)

func AclLoad(){
	var tmp []ACLRole
	cur, err := mongohelper.Database.Collection("acl").Find(context.TODO(), bson.D{})
	iohelper.ErrLog(err)
	err = cur.All(context.TODO(), &tmp)
	iohelper.ErrLog(err)
	for _, e := range tmp{
		AclRoles[e.Name] = e.Include
	}
}

func AclAllow(a string, b string) bool{
	if a == b{
		return true
	}
	for _, e := range AclRoles[a]{
		if e == b{
			return true
		}
	}
	return false
}
