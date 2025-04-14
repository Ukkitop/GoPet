package middleware

import (
	"errors"
	"net/http"
	"realAPI/api/apiUtils"
	"strings"
)

var UnAuthorizedError = errors.New("unauthorized")

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token = r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		err := apiUtils.VerifyToken(token)

		if err != nil {
			apiUtils.RequestErrorHandler(w, UnAuthorizedError)
		}

		next.ServeHTTP(w, r)
	})
}
