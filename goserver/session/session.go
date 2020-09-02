package session

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var sessionCookieSecret = flag.String("session_cookie_secret", "", "secret key used to encrypt the session cookie")

const (
	sessionLifeInSec = 7200
	cookieLifeInSec  = 7200
	cookieName       = "session"
)

// Session stores necessary information for a telehealth session.
type Session struct {
	ID, FHIRURL string
	expireAt    time.Time
}

// SessionID the id of the session.
func (s *Session) SessionID() string {
	return s.sid
}

// New creates a new session and set cookie containning the encoded session id.
func New(m Manager, w http.ResponseWriter, r *http.Request) (*Session, error) {
	expireAt := time.Now().Add(sessionLifeInSec * time.Second)
	s, err := m.Create(expireAt)
	if err != nil {
		return nil, err
	}
	expiration := time.Now().Add(cookieLifeInSec * time.Second)
	cookie := &http.Cookie{Name: cookieName, Value: encodeSessionID(s.SessionID()), Expires: expiration}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	return s, nil
}

// Find returns the session in session manager matching the session id in the cookie of the request.
func Find(m Manager, r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	sid, err := decodeSessionID(cookie.Value)
	if err != nil {
		return nil, err
	}

	return m.Find(sid)
}

func encodeSessionID(sid string) string {
	b := base64.StdEncoding.EncodeToString([]byte(sid))
	s := fmt.Sprintf("%s-%s", b, signature(sid))
	return url.QueryEscape(s)
}

func signature(sid string) string {
	h := hmac.New(sha1.New, []byte(*sessionCookieSecret))
	h.Write([]byte(sid))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func decodeSessionID(esid string) (string, error) {
	esid, err := url.QueryUnescape(esid)
	if err != nil {
		return "", err
	}

	vals := strings.Split(esid, "-")
	if len(vals) != 2 {
		return "", fmt.Errorf("Invalid session ID")
	}

	bsid, err := base64.StdEncoding.DecodeString(vals[0])
	if err != nil {
		return "", err
	}
	sid := string(bsid)

	sig := signature(sid)
	if sig != vals[1] {
		return "", fmt.Errorf("Invalid session ID")
	}
	return sid, nil
}
