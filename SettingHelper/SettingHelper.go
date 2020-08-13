package SettingHelper

import (
	"github.com/devngho/openN-Go/IOHelper"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

//Folder and File Init
func InitFolderFile(){
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	IOHelper.ErrLog(err)

	_, err = os.Stat(filepath.Join(dir, "setting.ini"))

	//File Exist check
	if err != nil {
		if os.IsNotExist(err) {
			//Create ini File
			_ = IOHelper.CreateFile(filepath.Join(dir, "setting.ini"))
			cfg, err := ini.Load(filepath.Join(dir, "setting.ini"))
			IOHelper.ErrLog(err)
			_, _ = cfg.NewSection("cpu")
			cfg.Section("cpu").Key("core").MustInt(-1)
			_, _ = cfg.NewSection("default")
			cfg.Section("default").Key("namespace").MustString("문서")
			_, _ = cfg.NewSection("wiki")
			cfg.Section("wiki").Key("name").MustString("openN-Go 위키")
			_, _ = cfg.NewSection("wiki")
			cfg.Section("wiki").Key("license_html").MustString("<a rel=\"license\" href=\"http://creativecommons.org/licenses/by/4.0/\"><img alt=\"크리에이티브 커먼즈 라이선스\" style=\"border-width:0\" src=\"https://i.creativecommons.org/l/by/4.0/88x31.png\" /></a><br />이 저작물은 <a rel=\"license\" href=\"http://creativecommons.org/licenses/by/4.0/\">크리에이티브 커먼즈 저작자표시 4.0 국제 라이선스</a>에 따라 이용할 수 있습니다.")
			_, _ = cfg.NewSection("wiki")
			cfg.Section("wiki").Key("name_next").MustString("는")
			_, _ = cfg.NewSection("wiki")
			cfg.Section("wiki").Key("start_page").MustString("대문")
			_ = cfg.SaveTo(filepath.Join(dir, "setting.ini"))
		} else {
			log.Fatal(err)
			return
		}
	}

	//Folder Exist Check
	if _, err := os.Stat(filepath.Join(dir, "db")); os.IsNotExist(err) {
		//Create Folders
		IOHelper.CreateFolder(filepath.Join(dir, "db"), 777)
	}
	if _, err := os.Stat(filepath.Join(dir, "db", "old")); os.IsNotExist(err) {
		//Create Folders
		IOHelper.CreateFolder(filepath.Join(dir, "db", "old"), 777)
	}
	if _, err := os.Stat(filepath.Join(dir, "theme")); os.IsNotExist(err) {
		//Create Folders
		IOHelper.CreateFolder(filepath.Join(dir, "theme"), 777)
	}
}

//Read settings with section and key
func ReadSetting(section string, key string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	IOHelper.ErrLog(err)

	_, err = os.Stat(filepath.Join(dir, "setting.ini"))
	IOHelper.ErrLog(err)
	cfg, err := ini.Load(filepath.Join(dir, "setting.ini"))
	IOHelper.ErrLog(err)
	return cfg.Section(section).Key(key).String()
}