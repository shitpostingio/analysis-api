package structs

import "go.mongodb.org/mongo-driver/mongo/options"

// DatabaseConfiguration represents a MongoDB configuration.
type DatabaseConfiguration struct {
	DatabaseName   string
	Username       string
	Password       string
	AuthSource     string
	CollectionName string
	Hosts          []string
}

// MongoDBConnectionOptions represents the connection options for MongoDB.
func (c *DatabaseConfiguration) MongoDBConnectionOptions() *options.ClientOptions {

	clientOptions := options.Client()
	//clientOptions.SetAuth(options.Credential{
	//	AuthMechanism: "SCRAM-SHA-1",
	//	AuthSource:    c.AuthSource,
	//	Username:      c.Username,
	//	Password:      c.Password,
	//	PasswordSet:   true,
	//})

	clientOptions.SetHosts(c.Hosts)
	return clientOptions

}
