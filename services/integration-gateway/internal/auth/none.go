package auth

import "net/http"

// NoneAuthenticator sem autenticação (APIs públicas)
type NoneAuthenticator struct{}

func (a *NoneAuthenticator) Authenticate(req *http.Request) error {
	// Nada a fazer
	return nil
}

func (a *NoneAuthenticator) Type() string {
	return "none"
}
