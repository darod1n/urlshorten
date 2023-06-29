package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type logger interface {
	Errorf(template string, args ...interface{})
}

type KeyUserID string

func WithAutorization(h http.Handler, secretKey string, l logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := r.Cookie("Authorization-Token")
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				l.Errorf("failed to get Authorization-Token: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userID := uuid.NewString()
			token, err := buildJWTString(userID, secretKey)
			if err != nil {
				l.Errorf("failed to build jwt string: %v", err)
			}

			cookie := &http.Cookie{
				Name:  "Authorization-Token",
				Value: token,
			}
			http.SetCookie(w, cookie)

			ctx := context.WithValue(r.Context(), KeyUserID("UserID"), userID)
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := getUserID(authToken.Value, secretKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUserID("UserID"), userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(tokenString, secretKey string) (string, error) {
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

func buildJWTString(userID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to get signed string: %v", err)
	}

	return tokenString, nil
}
