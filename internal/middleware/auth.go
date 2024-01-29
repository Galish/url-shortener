// Package middleware implements common middlewares for http handlers.
package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/Galish/url-shortener/internal/auth"
	"github.com/Galish/url-shortener/internal/logger"
)

// Generic authentication middleware constants.
const (
	AuthCookieName = "auth"
	AuthHeaderName = "X-User"
)

var errMissingUserID = errors.New("user id not specified")

// WithAuthToken generates and stores a token in a cookie if one is not specified.
func WithAuthToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserID(r)
		if err != nil {
			logger.Debug(err)
		}

		isAuthorized := userID != ""
		if !isAuthorized {
			userID = uuid.NewString()
		}

		r.Header.Set(AuthHeaderName, userID)

		if isAuthorized {
			h.ServeHTTP(w, r)
			return
		}

		token, err := auth.NewToken(&auth.JWTClaims{
			UserID: userID,
		})
		if err != nil {
			logger.WithError(err).Error("unable to generate auth token")
			h.ServeHTTP(w, r)
			return
		}

		http.SetCookie(
			w,
			&http.Cookie{
				Name:  AuthCookieName,
				Value: token,
			},
		)

		h.ServeHTTP(w, r)
	})
}

// WithAuthChecker serves 401 error if authorization fails.
func WithAuthChecker(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserID(r)
		if err != nil {
			logger.Debug(err)
		}

		if err != nil {
			logger.WithError(err).Debug("unauthorized access attempt")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set(AuthHeaderName, userID)
		h.ServeHTTP(w, r)
	})
}

func getUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		return "", fmt.Errorf("unable to extract auth cookie: %w", err)
	}

	claims, err := auth.ParseToken(cookie.Value)
	if err != nil {
		return "", fmt.Errorf("unable to parse auth token: %w", err)
	}
	if claims.UserID == "" {
		return "", errMissingUserID
	}

	return claims.UserID, nil
}
