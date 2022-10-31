package authdomain

import "net/http"

type AuthReturn struct {
	// The redirect url.
	RedirectURL string
	// AuthToken is the auth token.
	AuthToken string
	//
	RespCode *http.Response
}

type CallbackRequest struct {
	// Cookie is the cookie.
	Cookie string
	// The state.
	State string
	// The code.
	Code string
	// The scope.
	Scope string
	// The authuser.
	AuthUser string
	// The prompt.
	Prompt string
}
