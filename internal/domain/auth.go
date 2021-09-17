package domain

import (
	"time"
)

type Session struct {
	Id           int64     `db:"id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	UserId       int64     `db:"user_id"`
}

func (s Session) Expired() bool {
	return s.ExpiresAt.Before(time.Now().UTC())
}

type SessionToUpdate struct {
	RefreshToken string
	ExpiresAt    time.Time
}

type Tokens struct {
	// Token used for accessing operations and/or resources
	AccessToken          string `json:"accessToken" binding:"required" example:"access token"`
	AccessTokenExpiredAt int16  `json:"-"`
	// Token used for refreshing session
	RefreshToken          string `json:"refreshToken" binding:"required" example:"refresh token"`
	RefreshTokenExpiredAt int16  `json:"-"`
} // @name Tokens
