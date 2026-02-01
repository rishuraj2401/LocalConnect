package models

import "time"

// Message represents a message between users
type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	Read       bool
	CreatedAt  time.Time
}

// Conversation represents a conversation between two users
type Conversation struct {
	ID               string
	User1ID          string
	User2ID          string
	LastMessageID    string
	LastMessageText  string
	LastMessageAt    time.Time
	UnreadCountUser1 int
	UnreadCountUser2 int
	CreatedAt        time.Time
}
