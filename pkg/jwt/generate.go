package jwt

import (
	"errors"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidAccessToken ...
	ErrInvalidAccessToken = errors.New("invalid access token")
)

// JWT ...
type JWT struct {
	RefreshSecret   []byte
	RefreshLifetime time.Duration
	AccessSecret    []byte
	AccessLifetime  time.Duration
	MailSecret      []byte
	MailLifetime    time.Duration
	PasswordCoast   int
}

// Token ...
type Token struct {
	AccessToken  string `json:"accessToken"`
	Refreshtoken string `json:"refreshToken"`
}

// Claims ...
type Claims struct {
	SessionID string
	jwtgo.StandardClaims
}

// New ...
func New(
	refreshSecret,
	accessSecret,
	mailSecret []byte,
	refreshLifetime,
	accessLifetime,
	mailLifetime time.Duration,
) *JWT {
	return &JWT{
		RefreshSecret:   refreshSecret,
		AccessSecret:    accessSecret,
		MailSecret:      mailSecret,
		RefreshLifetime: refreshLifetime,
		AccessLifetime:  accessLifetime,
		MailLifetime:    mailLifetime,
		PasswordCoast:   bcrypt.DefaultCost,
	}
}

// PasswordsMatch ...
func (jwt *JWT) PasswordsMatch(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Hash ...
func (jwt *JWT) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), jwt.PasswordCoast)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (jwt *JWT) newRefresh(sessionID string) (string, error) {
	return jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, Claims{
		SessionID: sessionID,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: time.Now().Add(jwt.RefreshLifetime * time.Minute).Unix(),
		},
	}).SignedString(jwt.RefreshSecret)
}

func (jwt *JWT) newAccess(sessionID string) (string, error) {
	return jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, Claims{
		SessionID: sessionID,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: time.Now().Add(jwt.AccessLifetime * time.Minute).Unix(),
		},
	}).SignedString(jwt.AccessSecret)
}

// GetClaimsAccess ...
func (jwt *JWT) GetClaimsAccess(access string) (*Claims, error) {
	claims := &Claims{}
	accToken, err := jwtgo.ParseWithClaims(access, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}
		return jwt.AccessSecret, nil
	})
	if err != nil {
		return claims, err
	}

	if !accToken.Valid {
		return claims, ErrInvalidAccessToken
	}

	return claims, nil
}

// GetClaimsRefresh ...
func (jwt *JWT) GetClaimsRefresh(refresh string) (*Claims, error) {
	claims := &Claims{}
	refToken, err := jwtgo.ParseWithClaims(refresh, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}
		return jwt.RefreshSecret, nil
	})
	if err != nil {
		return claims, err
	}

	if !refToken.Valid {
		return claims, ErrInvalidAccessToken
	}
	return claims, nil
}

// NewToken ...
func (jwt *JWT) NewToken(sessionID string) (Token, error) {
	token := Token{}
	var err error
	token.AccessToken, err = jwt.newAccess(sessionID)
	if err != nil {
		return token, err
	}

	token.Refreshtoken, err = jwt.newRefresh(sessionID)
	if err != nil {
		return token, err
	}
	return token, nil
}

// NewMailToken ...
func (jwt *JWT) NewMailToken(email string) (string, error) {
	return jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, Claims{
		SessionID: email,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: time.Now().Add(jwt.MailLifetime * time.Minute).Unix(),
		},
	}).SignedString(jwt.MailSecret)
}

// GetClaimsMail ...
func (jwt *JWT) GetClaimsMail(token string) (*Claims, error) {
	claims := &Claims{}
	mailToken, err := jwtgo.ParseWithClaims(token, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAccessToken
		}
		return jwt.MailSecret, nil
	})
	if err != nil {
		return claims, err
	}

	if !mailToken.Valid {
		return claims, ErrInvalidAccessToken
	}
	return claims, nil
}
