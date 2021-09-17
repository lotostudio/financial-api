package service

import (
	"context"
	"github.com/lotostudio/financial-api/internal/domain"
	"github.com/lotostudio/financial-api/internal/repo"
	"github.com/lotostudio/financial-api/pkg/auth"
	"github.com/lotostudio/financial-api/pkg/hash"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type AuthService struct {
	repo            repo.Users
	sessionsRepo    repo.Sessions
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func newAuthService(repo repo.Users, sessionsRepo repo.Sessions, hasher hash.PasswordHasher,
	tokenManager auth.TokenManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:            repo,
		sessionsRepo:    sessionsRepo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, toCreate domain.UserToCreate) (domain.User, error) {
	passwordHash, err := s.hasher.Hash(toCreate.Password)

	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		Email:     toCreate.Email,
		FirstName: toCreate.FirstName,
		LastName:  toCreate.LastName,
		Password:  passwordHash,
	}

	userId, err := s.repo.Create(ctx, user)

	if err != nil {
		return domain.User{}, err
	}

	// New user can be returned, so create session async
	go func() {
		if err := s.sessionsRepo.Create(context.Background(), userId); err != nil {
			log.Warnf("error creating empty session for user %d error - %s", userId, err)
		}
	}()

	return s.repo.Get(ctx, userId)
}

func (s *AuthService) Login(ctx context.Context, toLogin domain.UserToLogin) (domain.Tokens, error) {
	passwordHash, err := s.hasher.Hash(toLogin.Password)

	if err != nil {
		return domain.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, toLogin.Email, passwordHash)

	if err != nil {
		return domain.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *AuthService) Refresh(ctx context.Context, token string) (domain.Tokens, error) {
	session, err := s.sessionsRepo.GetByToken(ctx, token)

	if err != nil {
		return domain.Tokens{}, err
	}

	if session.Expired() {
		return domain.Tokens{}, ErrRefreshTokenExpired
	}

	return s.createSession(ctx, session.UserId)
}

func (s *AuthService) createSession(ctx context.Context, userId int64) (domain.Tokens, error) {
	var res domain.Tokens
	var err error

	res.AccessToken, err = s.tokenManager.Issue(strconv.FormatInt(userId, 10))

	if err != nil {
		return res, err
	}

	res.AccessTokenExpiredAt = int16(s.accessTokenTTL.Seconds())

	res.RefreshToken, err = s.tokenManager.Random()

	if err != nil {
		return res, nil
	}

	res.RefreshTokenExpiredAt = int16(s.refreshTokenTTL.Seconds())

	session := domain.SessionToUpdate{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().UTC().Add(s.refreshTokenTTL),
	}

	_, err = s.sessionsRepo.Update(ctx, session, userId)

	return res, err
}
