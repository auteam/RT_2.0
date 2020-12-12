package controller

import (
	"net/http"
	"time"

	"../model"
	"github.com/dgrijalva/jwt-go"
)

func (s *server) myjwt(w http.ResponseWriter, r *http.Request, Email, Password string) (*model.User, string, string, int) {

	u, err := s.store.User().FindByEmail(Email)

	if err != nil {
		s.error(w, r, http.StatusTeapot, errIncorrectEmail) //errIncorrectEmail)
		return nil, "", "", http.StatusUnauthorized
	}
	// if err := u.Validate(); err != nil {
	// 	s.error(w, r, http.StatusUnauthorized, errIncorrectEmail)
	// 	return nil, "", "", http.StatusUnauthorized
	// }
	if !u.ComparePassword(Password) {
		s.error(w, r, http.StatusUnauthorized, errIncorrectPassword)
		return nil, "", "", http.StatusUnauthorized
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	expirationTimeRefresh := time.Now().Add(5 * 60 * 24 * time.Minute)
	claimsRefresh := &Claims{
		Username: u.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeRefresh.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, "", "", http.StatusInternalServerError
	}

	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	tokenStringRefresh, err := tokenRefresh.SignedString([]byte(s.secretKeyRefresh))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, "", "", http.StatusInternalServerError
	}

	return u, tokenString, tokenStringRefresh, 0

}
