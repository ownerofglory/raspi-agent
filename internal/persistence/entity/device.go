package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Device represents a row in the `devices` table
type Device struct {
	ID               uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	Name             string    `gorm:"type:varchar(256);default:''"`
	OTP              string    `gorm:"type:varchar(256);default:''"`
	EnrollmentStatus string    `gorm:"type:varchar(256);default:''"`
	UserID           uuid.UUID `gorm:"type:uuid;"`
	User             *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// BeforeCreate hook to auto-generate UUIDs
func (d *Device) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID, err = uuid.NewV7()
		return
	}
	return
}
