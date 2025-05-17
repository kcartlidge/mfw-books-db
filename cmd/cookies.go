package main

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

// CookieHandler handles cookie operations
type CookieHandler struct {
	sc         *securecookie.SecureCookie
	altCookies bool
}

// NewCookieHandler creates a new cookie handler
func NewCookieHandler(altCookies bool) (*CookieHandler, error) {
	// Generate a random hash key for cookie encryption
	hashKey := securecookie.GenerateRandomKey(32)
	if hashKey == nil {
		return nil, nil
	}

	// Generate a random block key for cookie encryption
	blockKey := securecookie.GenerateRandomKey(32)
	if blockKey == nil {
		return nil, nil
	}

	return &CookieHandler{
		sc:         securecookie.New(hashKey, blockKey),
		altCookies: altCookies,
	}, nil
}

// SetCookie sets a secure cookie
func (ch *CookieHandler) SetCookie(w http.ResponseWriter, name string, value string, maxAge int) error {
	// Encode the value
	encoded, err := ch.sc.Encode(name, value)
	if err != nil {
		return err
	}

	// Create the cookie
	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,

		// NOT Safari
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// The following should NEVER be used when serving sites externally,
	// but are needed for local development because Safari will not
	// allow cookies to be set against 'localhost' without them.
	if ch.altCookies {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteNoneMode
	}

	// Set the cookie
	http.SetCookie(w, cookie)
	return nil
}

// GetCookie gets a secure cookie value
func (ch *CookieHandler) GetCookie(r *http.Request, name string) (string, error) {
	// Get the cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	// Decode the value
	var value string
	err = ch.sc.Decode(name, cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

// DeleteCookie deletes a cookie
func (ch *CookieHandler) DeleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
