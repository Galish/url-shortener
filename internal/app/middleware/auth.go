package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Galish/url-shortener/internal/app/auth"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

const (
	AuthCookieName = "auth"
	AuthHeaderName = "X-User"
)

var errMissingUserID = errors.New("user id not specified")

func WithAuthToken(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			h(w, r)
			return
		}

		token, err := auth.NewToken(&auth.JWTClaims{
			UserID: userID,
		})
		if err != nil {
			logger.WithError(err).Error("unable to generate auth token")
			h(w, r)
			return
		}

		http.SetCookie(
			w,
			&http.Cookie{
				Name:  AuthCookieName,
				Value: token,
			},
		)

		h(w, r)
	}
}

func WithAuthChecker(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserID(r)
		if err != nil {
			logger.Debug(err)
		}

		if errors.Is(err, errMissingUserID) {
			logger.Debug("unauthorized access attempt")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set(AuthHeaderName, userID)
		h(w, r)
	}
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
