package userhelper

import (
	"context"
	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/mongohelper"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/openpgp/errors"
	"os"
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

func SaveState() {
	_, err := mongohelper.Database.Collection("user").DeleteMany(context.TODO(), bson.D{})
	iohelper.ErrLog(err)
	b := make([]interface{}, len(Users))
	for i := range Users {
		b[i] = Users[i]
	}
	_, err = mongohelper.Database.Collection("user").InsertMany(context.TODO(), b)
	iohelper.ErrLog(err)
}

func Load()  {
	cur, err := mongohelper.Database.Collection("user").Find(context.TODO(), bson.D{})
	iohelper.ErrLog(err)
	err = cur.All(context.TODO(), &Users)
	iohelper.ErrLog(err)
}