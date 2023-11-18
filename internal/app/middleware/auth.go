package middleware

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/auth"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

const (
	AuthCookieName = "auth"
	AuthHeaderName = "X-User"
)

func WithAuthToken(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, isAuthorized := getUserDetails(r)

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
		userID, isAuthorized := getUserDetails(r)

		if isAuthorized {
			r.Header.Set(AuthHeaderName, userID)
			h(w, r)
			return
		}

		logger.Debug("unauthorized access attempt")

		r.Header.Del(AuthHeaderName)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func getUserDetails(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(AuthCookieName)
	if err != nil {
		logger.WithError(err).Debug("unable to extract auth cookie")
	}
	if cookie == nil {
		return uuid.New().String(), false
	}

	claims, err := auth.ParseToken(cookie.Value)
	if err != nil {
		logger.WithError(err).Debug("unable to parse auth token")
		return uuid.New().String(), false
	}
	if claims.UserID == "" {
		return uuid.New().String(), false
	}

	return claims.UserID, true
}
