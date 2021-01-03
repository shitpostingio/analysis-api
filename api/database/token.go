package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthorizedToken represents a token in the document store.
type AuthorizedToken struct {
	ID    primitive.ObjectID
	Token string
}

// IsAuthorized returns true if the token is found in the document store.
func IsAuthorized(token string) bool {

	if token == "" {
		return false
	}

	//
	ctx, cancelCtx := context.WithTimeout(dsCtx, opDeadline)
	defer cancelCtx()

	//
	filter := bson.M{"token": token}

	err := tokenCollection.FindOne(ctx, filter, options.FindOne()).Decode(&AuthorizedToken{})
	return err == nil

}
