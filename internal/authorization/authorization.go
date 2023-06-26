package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type Storage interface {
	CreateUserID(ctx context.Context) (string, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
}

func WithAutorization(h http.Handler, db Storage, secretKey string, l logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := r.Cookie("Authorization-Token")
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				l.Errorf("failed to get Authorization-Token: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userID, err := db.CreateUserID(r.Context())
			if err != nil {
				l.Errorf("failed to create UserID: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			token, err := BuildJWTString(userID, secretKey)
			if err != nil {
				l.Errorf("failed to build jwt string: %v", err)
			}

			cookie := &http.Cookie{
				Name:  "Authorization-Token",
				Value: token,
			}
			r.AddCookie(cookie)

			ctx := context.WithValue(r.Context(), "UserID", userID)
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := GetUserID(authToken.Value, secretKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "UserID", userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(tokenString, secretKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("not valid token: %v", token)
	}
	return claims.UserID, nil
}

func BuildJWTString(userID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to get signed string: %v", err)
	}

	return tokenString, nil
}
