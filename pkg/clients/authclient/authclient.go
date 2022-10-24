// Package authclient provides a client for the auth service.
package authclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/UArt-project/UArt-proxy/domain/authdomain"
)

type AuthClient interface {
	// SendAuthRequest sends an auth request.
	SendAuthRequest() (string, error)
	// SendOAuthData sends the OAuth data.
	SendOAuthData(callbackData authdomain.CallbackRequest) (*authdomain.AuthReturn, error)
}

// AuthServiceClient is a client for the auth service.
type AuthServiceClient struct {
	// The url of the auth service.
	url string
	// The timeout for the requests.
	timeout time.Duration
}

// NewAuthServiceClient creates a new instance of the AuthServiceClient.
func NewAuthServiceClient(url string, timeout time.Duration) *AuthServiceClient {
	return &AuthServiceClient{
		url:     url,
		timeout: timeout,
	}
}

// SendAuthRequest sends an auth request.
func (c AuthServiceClient) SendAuthRequest() (string, error) {
	// send auth request to the /auth endpoint and receive a 304 redirect to the auth service
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/oauth2/authorization/google", nil)
	if err != nil {
		return "", fmt.Errorf("creating request for sending an auth request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending an auth request: %w", err)
	}

	// get the redirect url
	redirectURL := resp.Header.Get("Location")

	err = resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("closing the response body: %w", err)
	}

	return redirectURL, nil
}

// SendOAuthData sends the OAuth data.
func (c AuthServiceClient) SendOAuthData(callbackData authdomain.CallbackRequest) (*authdomain.AuthReturn, error) {
	// send the OAuth data to the /auth/redirect endpoint and receive a 304 redirect to the proxy
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// getURL := fmt.Sprintf("%s%s?=%s", c.url, "/login/oauth2/code/google", url.QueryEscape(callbackData.State))
	getURL, err := url.Parse(c.url + "/login/oauth2/code/google")
	if err != nil {
		return nil, fmt.Errorf("parsing the url: %w", err)
	}

	query := getURL.Query()
	query.Add("state", callbackData.State)
	query.Add("code", callbackData.Code)
	query.Add("scope", callbackData.Scope)
	query.Add("authuser", callbackData.AuthUser)
	query.Add("prompt", callbackData.Prompt)

	getURL.RawQuery = query.Encode()

	fmt.Println(getURL.String())

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

	defer cancel()

	reqCode, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for sending the OAuth data: %w", err)
	}

	respCode, err := httpClient.Do(reqCode)
	if err != nil {
		return nil, fmt.Errorf("sending the OAuth data: %w", err)
	}

	// return &authdomain.AuthReturn{RespCode: respCode}, nil

	// get the redirect url
	// redirectURL := respCode.Header.Get("Location")

	// get session cookie
	sessionCookie := respCode.Header.Get("Set-Cookie")
	sessionCookie = strings.TrimRight(sessionCookie, "; Path=/; HttpOnly")

	// send get /id with the session cookie
	reqID, err := http.NewRequest("GET", c.url+"/id", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for sending the OAuth data: %w", err)
	}

	// reqID.AddCookie(sessionCookie)
	// if callbackData.Cookie != "" {
	// reqID.AddCookie(&http.Cookie{Value: callbackData.Cookie})
	// } else {
	reqID.Header.Set("Cookie", sessionCookie)
	// }

	respID, err := httpClient.Do(reqID)
	if err != nil {
		return nil, fmt.Errorf("sending the OAuth data: %w", err)
	}

	// fmt.Println("Request: ", reqID)

	// fmt.Println("Response: ", respID)

	// get auth token
	authToken := respID.Header.Get("Authorization")

	return &authdomain.AuthReturn{
		RedirectURL: "http://localhost:8888",
		AuthToken:   authToken,
	}, nil
}
