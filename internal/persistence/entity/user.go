package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a row in the `users` table
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"id"`
	FirstName    string    `gorm:"type:varchar(256);default:''" json:"first_name"`
	LastName     string    `gorm:"type:varchar(256);default:''" json:"last_name"`
	Email        string    `gorm:"type:varchar(256);not null;uniqueIndex" json:"email"`
	PasswordHash *string   `gorm:"type:text;" json:"-"`
	Provider     string    `gorm:"type:varchar(64);default:'local'" json:"provider"`
}

// BeforeCreate hook to auto-generate UUIDs
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID, err = uuid.NewV7()
		return
	}
	return
}
