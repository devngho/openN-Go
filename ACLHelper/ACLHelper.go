package ACLHelper

type ACL struct {
	Watch string `json:"watch"`
	Edit string `json:"edit"`
	AclEdit string `json:"acl_edit"`
}

type ACLRole struct {
	Name string `json:"name"`
	Include []string `json:"include"`
}

func (f ACL) AclAllow(a ACL){

}
