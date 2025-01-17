package db

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	mgm.DefaultModel `bson:",inline"`
	AccessToken      string             `json:"access_token" bson:"access_token"`
	User             primitive.ObjectID `json:"user" bson:"user"`
	RefreshToken     string             `json:"refresh_token" bson:"refresh_token"`
	Role             Role               `json:"role" bson:"role"`
	AccessExpiresAt  time.Time          `json:"access_expires_at" bson:"access_expires_at"`
	RefreshExpiresAt time.Time          `json:"refresh_expires_at" bson:"refresh_expires_at"`
	Blacklisted      bool               `json:"blacklisted" bson:"blacklisted"`
}

func NewToken(tokenAccess string, tokenRefresh string, role Role, expiresAt time.Time, refreshExpiresAt time.Time) *Token {
	return &Token{
		User:             primitive.NewObjectID(),
		AccessToken:      tokenAccess,
		RefreshToken:     tokenRefresh,
		AccessExpiresAt:  expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		Role:             role,
		Blacklisted:      false,
	}
}

func (model *Token) CollectionName() string {
	return "tokens"
}
