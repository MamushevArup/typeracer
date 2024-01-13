package mongodb

//
//import (
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//func NewClient(ctx context.Context, port, host, username, password, database, authDb string) (*mongo.Database, error) {
//	var mongoURL string
//	var isAuth bool
//	if username == "" && password == "" {
//		mongoURL = fmt.Sprintf("mongodb://%s:%s", host, port)
//	} else {
//		isAuth = true
//		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
//	}
//	cliOp := options.Client().ApplyURI(mongoURL)
//	if isAuth {
//		if authDb == "" {
//			authDb = database
//		}
//		cliOp.SetAuth(options.Credential{
//			AuthSource: authDb,
//			Username:   username,
//			Password:   password,
//		})
//	}
//	client, err := mongo.Connect(ctx, cliOp)
//	if err != nil {
//		return nil, fmt.Errorf("error with connecting to the db error is : %v", err)
//	}
//	if err = client.Ping(ctx, nil); err != nil {
//		return nil, fmt.Errorf("error with ping the db error  is : %v", err)
//	}
//	return client.Database(database), nil
//}
