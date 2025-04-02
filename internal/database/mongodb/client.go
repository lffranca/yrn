package mongodb

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/yrn-go/yrn/pkg/yctx"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOption "go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"os"
	"strings"
)

const (
	EnvMongoUrl      = "MONGO_URL"
	EnvMongoDatabase = "MONGO_DATABASE"
)

func GetCollection(ctx *yctx.Context, collectionName string) (collection *mongo.Collection, err error) {
	var (
		mongoDatabase *mongo.Database
	)

	mongoDatabase, err = getDatabase(ctx)
	if err != nil {
		return
	}

	return mongoDatabase.Collection(collectionName), nil
}

func getDatabase(ctx *yctx.Context) (database *mongo.Database, err error) {

	databaseName := os.Getenv(EnvMongoDatabase)
	if databaseName == "" {
		return nil, errors.New("missing environment variable: " + EnvMongoDatabase)
	}

	var mongoClient *mongo.Client
	mongoClient, err = getClient(ctx)
	if err != nil {
		return
	}

	return mongoClient.Database(databaseName), nil
}

func getClient(ctx *yctx.Context) (mongoClient *mongo.Client, err error) {
	var (
		mongoURI      *string
		tlsConfigData *tls.Config
	)

	mongoURI, tlsConfigData, err = getMongoURI()
	if err != nil {
		return
	}

	mongoOptions := mongoOption.Client().
		ApplyURI(*mongoURI)

	if tlsConfigData != nil {
		mongoOptions.SetTLSConfig(tlsConfigData)
	}

	mongoClient, err = mongo.Connect(ctx.Context(), mongoOptions)
	if err != nil {
		return
	}

	return mongoClient, nil
}

func getMongoURI() (mongoURI *string, tlsConfigData *tls.Config, err error) {
	var (
		connString = os.Getenv(EnvMongoUrl)
		urlParsed  *url.URL
	)

	urlParsed, err = url.Parse(connString)
	if err != nil {
		return
	}

	queryStrings := urlParsed.Query()

	if strings.ToLower(queryStrings.Get("ssl")) == "true" {
		var certs []byte

		certs, err = os.ReadFile(queryStrings.Get("ssl_ca_certs"))
		if err != nil {
			return
		}

		tlsConfig := new(tls.Config)
		tlsConfig.RootCAs = x509.NewCertPool()
		if ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs); !ok {
			return
		}

		tlsConfigData = tlsConfig

		queryStrings.Del("ssl_ca_certs")
	}

	urlParsed.RawQuery = queryStrings.Encode()

	finalURI := urlParsed.String()

	return &finalURI, tlsConfigData, nil
}
