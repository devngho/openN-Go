package mongohelper

import (
	"context"
	"github.com/devngho/openN-Go/iohelper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect(ip string, database string){
	clientOptions := options.Client().ApplyURI(ip)
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	iohelper.ErrLog(err)
	err = Client.Ping(context.TODO(), nil)
	iohelper.ErrLog(err)
	Database = Client.Database(database)
}

func Disconnect()  {
	err := Client.Disconnect(context.TODO())
	iohelper.ErrLog(err)
}