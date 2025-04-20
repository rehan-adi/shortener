package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortKey struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Key       string             `bson:"key"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"createdAt"`
}

const (
	Available = "available"
	Used      = "used"
)
