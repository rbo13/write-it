package jwtservice

import "github.com/go-chi/jwtauth"

// JWT ...
type JWT struct {
	TokenAuth *jwtauth.JWTAuth
}

const jwtSecret = "5f7532af1ee4524945250f694b5bd06f44f9127bfc35924c457dfa7f68356798319d2d2c4bdce5aaee390cdc731585285e1e374fc1a88dcdbe3f21320b602aba"

// New ...
func New() *JWT {
	return &JWT{
		TokenAuth: jwtauth.New("HS256", []byte(jwtSecret), nil),
	}
}
