package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	Base
	Email           string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash    *string    `gorm:"size:255" json:"-"`
	Name            string     `gorm:"size:255" json:"name"`
	AvatarURL       *string    `gorm:"size:500" json:"avatar_url,omitempty"`
	Role            UserRole   `gorm:"type:varchar(20);default:user" json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`

	// Relations
	OAuthAccounts []OAuthAccount `gorm:"foreignKey:UserID" json:"oauth_accounts,omitempty"`
	Scans         []Scan         `gorm:"foreignKey:UserID" json:"scans,omitempty"`
	Corrections   []Correction   `gorm:"foreignKey:UserID" json:"corrections,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) HasPassword() bool {
	return u.PasswordHash != nil && *u.PasswordHash != ""
}

// SetPassword hashes and sets the user password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashStr := string(hash)
	u.PasswordHash = &hashStr
	return nil
}

// CheckPassword verifies a password against the stored hash
func (u *User) CheckPassword(password string) bool {
	if u.PasswordHash == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password))
	return err == nil
}

type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	Role            UserRole   `json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		Name:            u.Name,
		AvatarURL:       u.AvatarURL,
		Role:            u.Role,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
	}
}
