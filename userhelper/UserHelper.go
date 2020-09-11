package userhelper

import (
	"encoding/json"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/openpgp/errors"
	"os"
	"path/filepath"
)

type User struct {
	Acl            string `json:"acl"`
	Name           string `json:"name"`
	PasswordHashed [64]byte `json:"password_hashed"`
	Uid            string `json:"uid"`
}

var Users []User

func FindUserWithUid(Uid string) (User, error){
	for _, user := range Users {
		if user.Uid == Uid{
			return user, nil
		}
	}
	return User{}, os.ErrNotExist
}

func FindUserWithNamePwd(Name string, PwdHashed [64]byte) (User, error){
	for _, user := range Users {
		if user.Name == Name{
			if user.PasswordHashed == PwdHashed{
				return user, nil
			}else{
				return User{}, errors.ErrKeyIncorrect
			}
		}
	}
	return User{}, os.ErrNotExist
}

func Signup(Name string, PasswordHashed [64]byte)  {
	user := User{
		Acl:            "user",
		Name:           Name,
		PasswordHashed: PasswordHashed,
		Uid: ksuid.New().String(),
	}
	Users = append(Users, user)
	SaveState()
}

func SaveState()  {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	iohelper.ErrLog(err)
	f, err := os.Open(filepath.Join(dir, "db", "user.json"))
	if os.IsNotExist(err){

	}
	enc := json.NewEncoder(f)
	err = enc.Encode(&Users)
	iohelper.ErrLog(err)
	_ = f.Close()
}

func Load()  {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	iohelper.ErrLog(err)
	f, err := os.Open(filepath.Join(dir, "db", "user.json"))
	if os.IsNotExist(err){

	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&Users)
	iohelper.ErrLog(err)
	_ = f.Close()
}