package middleware

import (
	"net/http"

	"github.com/Galish/url-shortener/internal/app/auth"
	"github.com/Galish/url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

const authCookieName = "auth"

func WithAuthentication(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, isNewUser := getUserDetails(r)

		r.Header.Set("UserID", userID)

		if !isNewUser {
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
				Name:  authCookieName,
				Value: token,
			},
		)

		h(w, r)
	}
}

func getUserDetails(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		logger.WithError(err).Error("unable to extract auth cookie")
	}
	if cookie == nil {
		return uuid.New().String(), true
	}

	claims, err := auth.ParseToken(cookie.Value)
	if err != nil {
		logger.WithError(err).Error("unable to parse auth token")
		return uuid.New().String(), true
	}
	if claims.UserID == "" {
		return uuid.New().String(), true
	}

	return claims.UserID, false
}
