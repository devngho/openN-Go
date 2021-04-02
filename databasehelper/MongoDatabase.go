package databasehelper

import (
	"context"

	"github.com/devngho/openN-Go/iohelper"
	"github.com/devngho/openN-Go/types"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/sha3"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

type Mongodb struct{}

func init() {
	register("mongodb", Mongodb{})
}

func (t Mongodb) Connect(connection string, setting string) error {
	clientOptions := options.Client().ApplyURI(connection)
	var err error
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	err = MongoClient.Ping(context.TODO(), nil)
	MongoDatabase = MongoClient.Database(setting)
	return err
}

func (t Mongodb) CreateDocument(namespace string, name string, creator string) error {
	_, err := MongoDatabase.Collection("document").InsertOne(context.TODO(), types.Document{Namespace: namespace, Name: name, Editor: creator, Version: 1, Text: "", Acl: types.ACL{
		UseNamespace: true,
	}})
	return err
}

func (t Mongodb) CreateDocumentArchive(namespace string, name string, creator string) error {
	_, err := MongoDatabase.Collection("document_archive").InsertOne(context.TODO(), bson.M{"namespace": namespace, "name": name, "data": []types.DocumentArchive{{Namespace: namespace, Name: name, Editor: creator, Version: 1, Action: types.DocumentCreateAction, Text: "", Acl: types.ACL{
		UseNamespace: true,
	}}}})
	return err
}

func (t Mongodb) EditDocument(document types.Document, text string) error {
	//TODO
	return nil
}

func (t Mongodb) ReadDocument(namespace string, name string) (types.Document, error) {
	var u types.Document
	res := MongoDatabase.Collection("document").FindOne(context.TODO(), bson.M{"namespace": namespace, "name": name})
	if res.Err() != nil {
		return types.Document{}, res.Err()
	}
	err := res.Decode(&u)
	if err != nil {
		return types.Document{}, err
	}
	return u, err
}

func (t Mongodb) DeleteDocument(document types.Document) error {
	//TODO
	return nil
}

func (t Mongodb) AclEditDocument(document types.Document, acl types.ACL) error {
	//TODO
	return nil
}

func (t Mongodb) HasDocument(namespace string, name string) (bool, error) {
	count, err := MongoDatabase.Collection("document").CountDocuments(context.TODO(), bson.D{{"namespace", namespace}, {"name", name}})
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (t Mongodb) ReadACLRole() ([]types.ACLRole, error) {
	var tmp []types.ACLRole
	cur, err := MongoDatabase.Collection("acl").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &tmp)
	return tmp, err
}

func (t Mongodb) ReadNamespaces() ([]types.Namespace, error) {
	var results []types.Namespace // Use your own types here, but this works too

	cur, err := MongoDatabase.Collection("namespace").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &results)
	return results, err
}

func (t Mongodb) ReadUsers() ([]types.User, error) {
	var users []types.User
	cur, err := MongoDatabase.Collection("user").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cur.All(context.TODO(), &users)
	return users, err
}

func (t Mongodb) SaveUsers(user []types.User) error {
	_, err := MongoDatabase.Collection("user").DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	b := make([]interface{}, len(user))
	for i := range user {
		b[i] = user[i]
	}
	_, err = MongoDatabase.Collection("user").InsertMany(context.TODO(), b)
	return err
}

func (t Mongodb) InitData() error {
	result := MongoDatabase.Collection("user").FindOne(context.TODO(), bson.M{"name": "admin"})
	if result.Err() != nil {
		_, err := MongoDatabase.Collection("user").InsertOne(context.TODO(), types.User{Acl: "admin", PasswordHashed: sha3.Sum512([]byte("openngo")), Name: "admin", Uid: ksuid.New().String()})
		if err != nil {
			return err
		}
		_, err = MongoDatabase.Collection("acl").InsertMany(context.TODO(), []interface{}{types.ACLRole{Name: "admin", Include: []string{"ip", "user"}}, types.ACLRole{Name: "user", Include: []string{"ip"}}, types.ACLRole{Name: "ip", Include: []string{}}})
		iohelper.ErrLog(err)
		if err != nil {
			return err
		}
		_, err = MongoDatabase.Collection("namespace").InsertOne(context.TODO(), types.Namespace{Name: "문서", NamespaceACL: types.ACLNamespace{AclEdit: "admin", Edit: "ip", Watch: "ip", Create: "ip", Delete: "user"}})
		iohelper.ErrLog(err)
		if err != nil {
			return err
		}
	} else {
		return result.Err()
	}
	return nil
}
