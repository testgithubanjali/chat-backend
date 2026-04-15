package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderID   string             `bson:"sender_id"`
	ReceiverID string             `bson:"receiver_id"`
	Content    string             `bson:"content"`
	CreatedAt  time.Time          `bson:"created_at"`
}
