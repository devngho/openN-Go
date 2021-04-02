package types

import "time"

//Document

type Document struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Acl       ACL    `json:"acl"`
	Editor    string `json:"editor"`
	Version   int    `json:"version"`
}

type DocumentArchive struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Acl       ACL    `json:"acl"`
	Editor    string `json:"editor"`
	Version   int    `json:"version"`
	Action    string `json:"action"`
}

const (
	DocumentEditAction    = "edit"
	DocumentAclEditAction = "acl_edit"
	DocumentCreateAction  = "create"
	DocumentDeleteAction  = "delete"
)

//ACL

type ACL struct {
	Watch        string `json:"watch"`
	Edit         string `json:"edit"`
	AclEdit      string `json:"acl_edit"`
	Delete       string `json:"delete"`
	UseNamespace bool   `json:"use_namespace"`
}

type ACLNamespace struct {
	Watch        string `json:"watch"`
	Edit         string `json:"edit"`
	AclEdit      string `json:"acl_edit"`
	Create       string `json:"create"`
	Delete       string `json:"delete"`
	UseNamespace bool   `json:"use_namespace"`
}

type ACLRole struct {
	Name    string   `json:"name"`
	Include []string `json:"include"`
}

//User

type User struct {
	Acl            string   `json:"acl"`
	Name           string   `json:"name"`
	PasswordHashed [64]byte `json:"password_hashed"`
	Uid            string   `json:"uid"`
}

//Namespace

type Namespace struct {
	Name         string       `json:"name"`
	NamespaceACL ACLNamespace `json:"acl"`
}

//Markdown Parser

type MarkdownParser interface {
	ToHTML(markdown string) string
	ToMarkdown(html string) string
}

//Router Action

const WatchDocument = 0
const EditDocument = 1
const WatchArchiveListDocument = 2
const WatchRevDocument = 3
const Search = 4
const NewDocument = 5
const DeleteDocument = 6
const EditAclDocument = 7
const LoginPost = 8
const ExpirationTime = time.Duration(30) * time.Minute

//Setting Caching

var Screct string
var ScrectByte []byte
