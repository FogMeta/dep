package model

import "time"

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Email     string    `gorm:"column:email;size:128;default:''"`
	Password  string    `gorm:"column:password;size:64;default:''"`
	Wallet    string    `gorm:"column:wallet;size:128;default:''"`
	Type      int       `gorm:"column:type;type:tinyint;not null;default:0"`
	APIKey    string    `gorm:"column:api_key;size:64;default:''"`
	Status    int       `gorm:"column:status;type:tinyint;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "user"
}
