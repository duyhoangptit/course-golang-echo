package model

import "time"

type User struct {
	UserId    string    `json:"-" db:"user_id,omitempty"`
	FullName  string    `json:"full_name,omitempty" db:"full_name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Password  string    `json:"-" db:"password"`
	Role      string    `json:"-" db:"role"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	Token     string    `json:"token"`
}
