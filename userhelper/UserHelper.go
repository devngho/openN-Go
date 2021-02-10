package userhelper

import (
	"github.com/devngho/openN-Go/databasehelper"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/types"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/openpgp/errors"
	"os"
)

var Users []types.User

func FindUserWithUid(Uid string) (types.User, error) {
	for _, user := range Users {
		if user.Uid == Uid {
			return user, nil
		}
	}
	return types.User{}, os.ErrNotExist
}

func FindUserWithNamePwd(Name string, PwdHashed [64]byte) (types.User, error) {
	for _, user := range Users {
		if user.Name == Name {
			if user.PasswordHashed == PwdHashed {
				return user, nil
			} else {
				return types.User{}, errors.ErrKeyIncorrect
			}
		}
	}
	return types.User{}, os.ErrNotExist
}

func Signup(Name string, PasswordHashed [64]byte) {
	user := types.User{
		Acl:            "user",
		Name:           Name,
		PasswordHashed: PasswordHashed,
		Uid:            ksuid.New().String(),
	}
	Users = append(Users, user)
	SaveState()
}

func SaveState() {
	err := databasehelper.Dao.SaveUsers(Users)
	iohelper.ErrFatal(err)
}

func Load() {
	var err error
	Users, err = databasehelper.Dao.ReadUsers()
	iohelper.ErrFatal(err)
}
