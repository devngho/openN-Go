package ThemeHelper

import (
	"../SettingHelper"
	"strings"
	"io/ioutil"
)

var wikiName = SettingHelper.ReadSetting("wiki", "name")
var licenseHtml = SettingHelper.ReadSetting("wiki", "license_html")
var (
	documentHtml = ""
	documentAclBlockHtml = ""
	documentEditHtml = ""
	documentEditBlockHtml = ""
	documentOldHtml = ""
	emailCheckHtml = ""
	errorHtml = ""
	license = ""
)

func applyString(html string) string{
	html = strings.ReplaceAll(html, "${wiki}", wikiName)
	html = strings.ReplaceAll(html, "${license}", licenseHtml)
	return html
}

func ThemeDocument(namespace string, name string)  {

}