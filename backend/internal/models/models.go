package models

import "time"

type UserRole string

const (
	RoleWorker UserRole = "worker"
	RoleClient UserRole = "client"
)

type User struct {
	ID           string
	Name         string
	Email        string
	Phone        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
}

type Category struct {
	ID   int
	Name string
}

type WorkerProfile struct {
	ID              string
	UserID          string
	CategoryID      int
	CategoryName    string
	Location        string
	Rate            float64
	ExperienceYears int
	Bio             string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Upvotes         int
	ReviewCount     int
	AverageRating   float64
}

type Media struct {
	ID        string
	ProfileID string
	MediaType string
	URL       string
	CreatedAt time.Time
}

type Review struct {
	ID        string
	ProfileID string
	UserID    string
	Rating    int
	Comment   string
	CreatedAt time.Time
}

type ContactRequest struct {
	ID        string
	ProfileID string
	UserID    string
	Message   string
	Phone     string
	CreatedAt time.Time
}
