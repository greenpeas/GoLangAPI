package transport

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenNotValid = errors.New("Token is not valid")
)

type JwtWorker struct {
	Secret     string
	Ttl        time.Duration
	TtlRefresh time.Duration
}

type JwtWithRefresh struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type jwtClaim struct {
	UserId         int    `json:"user_id,omitempty"`
	UserRole       int    `json:"user_role,omitempty"`
	UserTitle      string `json:"user_title,omitempty"`
	IsRefreshToken bool   `json:"is_refresh_token,omitempty"`
	jwt.RegisteredClaims
}

func NewJwtWorker(secret string, ttl, ttlRefresh int64) JwtWorker {
	return JwtWorker{
		secret,
		time.Duration(ttl),
		time.Duration(ttlRefresh),
	}
}

func (w *JwtWorker) GenerateTokens(userId, userRole int, userTitle string) (JwtWithRefresh, error) {
	var jwrWithRefresh JwtWithRefresh
	var err error

	if jwrWithRefresh.Token, err = w.GenerateToken(userId, userRole, userTitle); err != nil {
		return JwtWithRefresh{}, err
	}

	if jwrWithRefresh.RefreshToken, err = w.GenerateRefreshToken(userId, userRole, userTitle); err != nil {
		return JwtWithRefresh{}, err
	}

	return jwrWithRefresh, nil
}

func (w *JwtWorker) GenerateToken(userId, userRole int, userTitle string) (string, error) {
	return w.generateToken(userId, userRole, userTitle, false)
}

func (w *JwtWorker) GenerateRefreshToken(userId, userRole int, userTitle string) (string, error) {
	return w.generateToken(userId, userRole, userTitle, true)
}

func (w *JwtWorker) generateToken(userId, userRole int, userTitle string, isRefreshToken bool) (string, error) {
	ttl := w.Ttl
	if isRefreshToken {
		ttl = w.TtlRefresh
	}

	claims := jwtClaim{
		UserId:         userId,
		UserRole:       userRole,
		UserTitle:      userTitle,
		IsRefreshToken: isRefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(w.Secret))
}

func (w *JwtWorker) TokenIsValid(accessToken string) bool {
	if claims, err := w.getValidClaimsFromToken(accessToken); err != nil {
		return false
	} else {
		return !claims.IsRefreshToken
	}
}

func (w *JwtWorker) RefreshTokenIsValid(accessToken string) bool {
	if claims, err := w.getValidClaimsFromToken(accessToken); err != nil {
		return false
	} else {
		return claims.IsRefreshToken
	}
}

func (w *JwtWorker) getValidClaimsFromToken(accessToken string) (*jwtClaim, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(w.Secret), nil
	})

	if err != nil {
		return &jwtClaim{}, err
	}

	if claims, ok := token.Claims.(*jwtClaim); ok && token.Valid {
		return claims, nil
	}

	return &jwtClaim{}, ErrTokenNotValid
}

func (w *JwtWorker) GetUserIdFromToken(accessToken string) (int, error) {
	if claims, err := w.getValidClaimsFromToken(accessToken); err != nil {
		return 0, err
	} else if claims.IsRefreshToken {
		return claims.UserId, ErrTokenNotValid
	} else {
		return claims.UserId, nil
	}
}

func (w *JwtWorker) GetUserIdFromRefreshToken(accessToken string) (int, error) {
	if claims, err := w.getValidClaimsFromToken(accessToken); err != nil {
		return 0, err
	} else if claims.IsRefreshToken {
		return claims.UserId, nil
	} else {
		return claims.UserId, ErrTokenNotValid
	}
}
