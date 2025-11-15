package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UserWithDevice(db *gorm.DB) error {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202511152058",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					ID           uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
					FirstName    string    `gorm:"type:varchar(256);default:''"`
					LastName     string    `gorm:"type:varchar(256);default:''"`
					Email        string    `gorm:"type:varchar(256);not null;uniqueIndex"`
					PasswordHash *string   `gorm:"type:text;" json:"-"`
					Provider     string    `gorm:"type:varchar(64);default:'local'"`
				}

				type Device struct {
					ID               uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
					Name             string    `gorm:"type:varchar(256);default:''"`
					OTP              string    `gorm:"type:varchar(256);default:''"`
					EnrollmentStatus string    `gorm:"type:varchar(256);default:''"`
					UserID           uuid.UUID `gorm:"type:uuid;"`
					User             *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
				}

				return tx.AutoMigrate(&User{}, &Device{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("devices", "users")
			},
		},
	}).Migrate()
}
