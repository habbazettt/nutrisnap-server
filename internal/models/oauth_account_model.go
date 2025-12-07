package models

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider string

const (
	ProviderGoogle OAuthProvider = "google"
	ProviderApple  OAuthProvider = "apple"
)

type OAuthAccount struct {
	BaseWithoutSoftDelete
	UserID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider       OAuthProvider `gorm:"type:varchar(50);not null" json:"provider"`
	ProviderUserID string        `gorm:"size:255;not null" json:"provider_user_id"`
	AccessToken    *string       `gorm:"size:2000" json:"-"`
	RefreshToken   *string       `gorm:"size:2000" json:"-"`
	TokenExpiresAt *time.Time    `json:"token_expires_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}

func (o *OAuthAccount) IsTokenExpired() bool {
	if o.TokenExpiresAt == nil {
		return true
	}
	return time.Now().After(*o.TokenExpiresAt)
}
