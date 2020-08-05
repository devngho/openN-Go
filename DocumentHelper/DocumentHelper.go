package DocumentHelper

import (
	"../IOHelper"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type ACL struct {
	Watch string
	Edit string
	AclEdit string
}
type Document struct {
	Namespace string `json:"namespace"`
	Name string `json:"name"`
	Text string `json:"text"`
	Acl ACL `json:"acl"`
	Editor string `json:"editor"`
	Version int `json:"version"`
}
func (d Document) AclAllow(a ACL){

}
func (d Document) Edit(n string){
	
}
func (d Document) Delete() {

}
func Read(Namespace string, Name string) (Document, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(filepath.Join(dir, "db", Namespace+"_"+Name+".json"))
	if os.IsNotExist(err){
		return Document{} ,err
	}
	dec := json.NewDecoder(f)
	u := Document{}
	err = dec.Decode(&u)
	IOHelper.ErrLog(err)
	return u, err
}