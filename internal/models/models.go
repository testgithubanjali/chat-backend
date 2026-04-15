package models
import "time"

type Message struct {
	ID   uint  `gorm:"primaryKey" json:"id"`
	SenderID string  `json:"sender_id"`
	ReceiverID string  `json:"receiver_id"`
	Content string  `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	
}