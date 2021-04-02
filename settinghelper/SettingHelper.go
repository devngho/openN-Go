package settinghelper

import (
	"log"
	"os"
	"path/filepath"

	"github.com/devngho/openN-Go/databasehelper"
	"github.com/devngho/openN-Go/iohelper"
	"gopkg.in/ini.v1"
)

var Setting *ini.File

func LoadSettings() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	iohelper.ErrLog(err)

	_, err = os.Stat(filepath.Join(dir, "setting.ini"))
	iohelper.ErrLog(err)
	Setting, err = ini.Load(filepath.Join(dir, "setting.ini"))
	iohelper.ErrLog(err)
}

//Folder and File Init
func InitFolderFile() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	iohelper.ErrLog(err)

	_, err = os.Stat(filepath.Join(dir, "setting.ini"))

	//File Exist check
	if err != nil {
		if os.IsNotExist(err) {
			//Create ini File
			_ = iohelper.CreateFile(filepath.Join(dir, "setting.ini"))
			cfg, err := ini.Load(filepath.Join(dir, "setting.ini"))
			iohelper.ErrLog(err)
			_, _ = cfg.NewSection("cpu")
			cfg.Section("cpu").Key("core").MustInt(-1)
			_, _ = cfg.NewSection("default")
			cfg.Section("default").Key("namespace").MustString("문서")
			_, _ = cfg.NewSection("wiki")
			cfg.Section("wiki").Key("name").MustString("openN-Go 위키")
			cfg.Section("wiki").Key("license_html").MustString("<a rel=\"license\" href=\"http://creativecommons.org/licenses/by/4.0/\"><img alt=\"크리에이티브 커먼즈 라이선스\" style=\"border-width:0\" src=\"https://i.creativecommons.org/l/by/4.0/88x31.png\" /></a><br />이 저작물은 <a rel=\"license\" href=\"http://creativecommons.org/licenses/by/4.0/\">크리에이티브 커먼즈 저작자표시 4.0 국제 라이선스</a>에 따라 이용할 수 있습니다.")
			cfg.Section("wiki").Key("name_next").MustString("는")
			cfg.Section("wiki").Key("start_page").MustString("대문")
			cfg.Section("wiki").Key("use_markdown").MustBool(true)
			cfg.Section("wiki").Key("markdown").MustString("example")
			cfg.Section("wiki").Key("use_only_api").MustBool(false)
			_, _ = cfg.NewSection("db")
			cfg.Section("db").Key("type").MustString("mongodb")
			cfg.Section("db").Key("server").MustString("mongodb://localhost:27017")
			cfg.Section("db").Key("setting").MustString("wiki")
			_, _ = cfg.NewSection("screct")
			cfg.Section("screct").Key("key").MustString("SCRECT")
			_ = cfg.SaveTo(filepath.Join(dir, "setting.ini"))
		} else {
			log.Fatal(err)
			return
		}
	}

	/*
		//Folder Exist Check
		if _, err := os.Stat(filepath.Join(dir, "db")); os.IsNotExist(err) {
			//Create Folder
			iohelper.CreateFolder(filepath.Join(dir, "db"), 777)
			//Create Namespace File
			f1, _ := os.Create(filepath.Join(dir, "db", "namespaces.json"))
			defer f1.Close()
			u := []namespacehelper.Namespace{{Name: "문서", NamespaceACL: aclhelper.ACLNamespace{AclEdit: "admin", Edit: "ip", Watch: "ip", Create: "ip", Delete: "user"}}}
			enc1 := json.NewEncoder(f1)
			enc1.SetIndent("", "  ")
			_ = enc1.Encode(u)
			//Create ACL File
			f2, _ := os.Create(filepath.Join(dir, "db", "acl.json"))
			defer f2.Close()
			ua := []aclhelper.ACLRole{{Name: "admin", Include: []string{"ip", "user"}}, {Name: "user", Include: []string{"ip"}}, {Name: "ip", Include: []string{}}}
			enc2 := json.NewEncoder(f2)
			enc2.SetIndent("", "  ")
			_ = enc2.Encode(ua)
			//Create User File
			f3, _ := os.Create(filepath.Join(dir, "db", "user.json"))
			defer f3.Close()
			uaa := []userhelper.User{{Acl: "admin", PasswordHashed: sha3.Sum512([]byte("openngo")), Name: "admin", Uid: ksuid.New().String()}}
			enc3 := json.NewEncoder(f3)
			enc3.SetIndent("", "  ")
			_ = enc3.Encode(uaa)
		}
		if _, err := os.Stat(filepath.Join(dir, "db", "old")); os.IsNotExist(err) {
			//Create Folder
			iohelper.CreateFolder(filepath.Join(dir, "db", "old"), 777)
		}
		if _, err := os.Stat(filepath.Join(dir, "theme")); os.IsNotExist(err) {
			//Create Folder
			iohelper.CreateFolder(filepath.Join(dir, "theme"), 777)
		}*/
}

func InitData() {
	iohelper.ErrFatal(databasehelper.Dao.InitData())
}

//Read settings with section and key
func ReadSetting(section string, key string) *ini.Key {
	return Setting.Section(section).Key(key)
}
