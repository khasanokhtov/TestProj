package models

import "time"

//Структур данных пользователя
type User struct{
	ID 		   uint `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	Balance    int       `json:"balance"`
	ReferrerID *uint     `json:"referrer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}