package db

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	mgm.DefaultModel `bson:",inline"`
	AccessToken      string             `json:"access_token" bson:"access_token"`             // Token truy cập
	User             primitive.ObjectID `json:"user" bson:"user"`                             // ID người dùng
	RefreshToken     string             `json:"refresh_token" bson:"refresh_token"`           // Token làm mới
	Role             Role               `json:"role" bson:"role"`                             // Vai trò người dùng
	AccessExpiresAt  time.Time          `json:"access_expires_at" bson:"access_expires_at"`   // Thời điểm hết hạn token truy cập
	RefreshExpiresAt time.Time          `json:"refresh_expires_at" bson:"refresh_expires_at"` // Thời điểm hết hạn token truy cập
	Blacklisted      bool               `json:"blacklisted" bson:"blacklisted"`               // Trạng thái blacklist
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

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
