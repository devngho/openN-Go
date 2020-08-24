package DocumentHelper

import (
	"encoding/json"
	"github.com/devngho/openN-Go/ACLHelper"
	"github.com/devngho/openN-Go/IOHelper"
	"os"
	"path/filepath"
)

type Document struct {
	Namespace string `json:"namespace"`
	Name string `json:"name"`
	Text string `json:"text"`
	Acl ACLHelper.ACL `json:"acl"`
	Editor string `json:"editor"`
	Version int `json:"version"`
}

func (d Document) Edit(n string){
	
}
func (d Document) Delete() {

}
func Create(namespace string, name string, creater string) (Document, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	IOHelper.ErrLog(err)
	f := IOHelper.CreateFile(filepath.Join(dir, "db", namespace+"_"+name+".json"))
	old := IOHelper.CreateFile(filepath.Join(dir, "db", "old", namespace+"_"+name+".json"))
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	u := Document{Namespace: namespace, Name: name, Editor: creater, Version: 1, Text: "", Acl: ACLHelper.ACL{
		Watch:   "ip",
		Edit:    "ip",
		AclEdit: "admin",
	}}
	_ = enc.Encode(u)
	_ = f.Close()
	enc = json.NewEncoder(old)
	enc.SetIndent("", "  ")
	uo := []Document{{Namespace: namespace, Name: name, Editor: creater, Version: 1, Text: "", Acl: ACLHelper.ACL{
		Watch:   "ip",
		Edit:    "ip",
		AclEdit: "admin",
	}}}
	_ = enc.Encode(uo)
	_ = f.Close()
	return u, err
}
func Read(Namespace string, Name string) (Document, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	IOHelper.ErrLog(err)
	f, err := os.Open(filepath.Join(dir, "db", Namespace+"_"+Name+".json"))
	if os.IsNotExist(err){
		return Document{} ,err
	}
	dec := json.NewDecoder(f)
	u := Document{}
	err = dec.Decode(&u)
	IOHelper.ErrLog(err)
	_ = f.Close()
	return u, err
}