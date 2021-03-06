package ThemeHelper

import (
	"../IOHelper"
	"../SettingHelper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var WikiName string
var WikiNext string
var LicenseHtml string
var (
	DocumentHtml = ""
	DocumentAclBlockHtml = ""
	DocumentEditHtml = ""
	DocumentEditBlockHtml = ""
	DocumentOldHtml = ""
	EmailCheckHtml = ""
	ErrorHtml = ""
	LicenseHtmlFile = ""
	LoginHtml = ""
	NotFoundDocumentHtml = ""
	NotFoundHtml = ""
	OverlapHtml = ""
	SignupHtml = ""
	WaitHtml = ""
)

func InitStatic()  {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	IOHelper.ErrLog(err)

	WikiName = SettingHelper.ReadSetting("wiki", "name")
	WikiNext = SettingHelper.ReadSetting("wiki", "name_next")
	LicenseHtml = SettingHelper.ReadSetting("wiki", "license_html")

	registryFileToVar(filepath.Join(dir, "theme", "document.html"), &DocumentHtml)
	registryFileToVar(filepath.Join(dir, "theme", "document_acl_block.html"), &DocumentAclBlockHtml)
	registryFileToVar(filepath.Join(dir, "theme", "document_edit.html"), &DocumentEditHtml)
	registryFileToVar(filepath.Join(dir, "theme", "document_edit_block.html"), &DocumentEditBlockHtml)
	registryFileToVar(filepath.Join(dir, "theme", "document_old.html"), &DocumentOldHtml)
	registryFileToVar(filepath.Join(dir, "theme", "email_check.html"), &EmailCheckHtml)
	registryFileToVar(filepath.Join(dir, "theme", "error.html"), &ErrorHtml)
	registryFileToVar(filepath.Join(dir, "theme", "license.html"), &LicenseHtmlFile)
	registryFileToVar(filepath.Join(dir, "theme", "login.html"), &LoginHtml)
	registryFileToVar(filepath.Join(dir, "theme", "notfound.html"), &NotFoundDocumentHtml)
	registryFileToVar(filepath.Join(dir, "theme", "notfound-not-document.html"), &NotFoundHtml)
	registryFileToVar(filepath.Join(dir, "theme", "overlap.html"), &OverlapHtml)
	registryFileToVar(filepath.Join(dir, "theme", "signup.html"), &SignupHtml)
	registryFileToVar(filepath.Join(dir, "theme", "wait.html"), &WaitHtml)
}

func registryFileToVar(file string, fileVar *string)  {
	read, err := ioutil.ReadFile(file)
	IOHelper.ErrLog(err)
	result := strings.ReplaceAll(string(read), "${wiki}", WikiName)
	result = strings.ReplaceAll(result, "${wikinext}", WikiNext)
	result = strings.ReplaceAll(result, "${license}", LicenseHtml)
	*fileVar = result
}