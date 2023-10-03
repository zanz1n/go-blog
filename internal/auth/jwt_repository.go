package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zanz1n/go-htmx/internal/errors"
	"github.com/zanz1n/go-htmx/internal/user"
	"github.com/zanz1n/go-htmx/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func NewJwtAuthService(hmac []byte, duration time.Duration) AuthRepository {
	return &JwtAuthRepository{
		hmacKey:       hmac,
		tokenDuration: duration,
	}
}

type JwtAuthRepository struct {
	tokenDuration time.Duration
	hmacKey       []byte
}

func (as *JwtAuthRepository) CreateUserToken(data *UserAuthPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, data)

	s, err := token.SignedString(as.hmacKey)
	if err != nil {
		return "", errors.ErrAuthTokenGenFailed
	}

	return s, nil
}

func (as *JwtAuthRepository) DecodeUserToken(payload string) (*UserAuthPayload, error) {
	claims := UserAuthPayload{}

	token, err := jwt.ParseWithClaims(payload, &claims, as.userTokenKeyFunc)
	if err != nil || !token.Valid {
		return nil, errors.ErrInvalidAuthToken
	}

	if err = claims.Validate(); err != nil {
		return nil, err
	}

	return &claims, nil
}

func (as *JwtAuthRepository) AuthUser(info *user.User, phash string) (string, error) {
	err := bcrypt.CompareHashAndPassword(
		utils.S2B(phash),
		utils.S2B(info.Password),
	)
	if err != nil {
		return "", errors.ErrLoginFailed
	}

	now := time.Now()
	claims := UserAuthPayload{
		UserId:     info.ID,
		Email:      info.Email,
		ExpiryDate: now.Add(as.tokenDuration).Unix(),
		IssuedAt:   now.Unix(),
		Role:       info.Role,
	}

	return as.CreateUserToken(&claims)
}

func (as *JwtAuthRepository) userTokenKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.ErrInvalidAuthToken
	}

	return as.hmacKey, nil
}

func (p *UserAuthPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(p.ExpiryDate, 0),
	}, nil
}

func (p *UserAuthPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: time.Unix(p.IssuedAt, 0),
	}, nil
}

func (p *UserAuthPayload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (p *UserAuthPayload) GetIssuer() (string, error) {
	return "", nil
}

func (p *UserAuthPayload) GetSubject() (string, error) {
	return p.UserId.String(), nil
}

func (p *UserAuthPayload) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}
