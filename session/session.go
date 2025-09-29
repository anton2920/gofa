package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type Session struct {
	/* ID is non-zero if user is authorized. */
	database.ID

	/* Customization is a user preferences for language, timezone, etc. */
	Customization

	/* Expiry is Unix time in seconds when session is no longer valid. */
	Expiry int64

	/* Token is key in Sessions map. */
	Token string
}

var (
	Sessions     = make(map[string]Session)
	SessionsLock sync.RWMutex
)

func New(userID database.ID) Session {
	token := GenerateToken()
	expiry := time.Unix() + time.OneWeek

	session := Session{ID: userID, Expiry: expiry, Token: token}
	SessionsLock.Lock()
	Sessions[token] = session
	SessionsLock.Unlock()

	return session
}

func Get(token string) Session {
	defer trace.End(trace.Begin(""))

	SessionsLock.RLock()
	session, ok := Sessions[token]
	SessionsLock.RUnlock()
	if !ok {
		return Session{}
	}

	now := time.Unix()
	if now-session.Expiry > 0 {
		SessionsLock.Lock()
		delete(Sessions, token)
		SessionsLock.Unlock()
		return Session{}
	}
	session.Expiry = now + time.OneWeek

	SessionsLock.Lock()
	Sessions[token] = session
	SessionsLock.Unlock()

	return session
}

func (session Session) RemoveAllForThisUser() {
	SessionsLock.Lock()
	for k, v := range Sessions {
		if v.ID == session.ID {
			delete(Sessions, k)
		}
	}
	SessionsLock.Unlock()
}

func (session Session) Update() {
	SessionsLock.Lock()
	Sessions[session.Token] = session
	SessionsLock.Unlock()
}

/* NOTE(anton2920): this function exists only so I don't have to solve name collisions of variable 'session' and this package. */
func (session Session) GenerateToken() string {
	return GenerateToken()
}

func GenerateToken() string {
	defer trace.End(trace.Begin(""))

	const n = 64
	buffer := make([]byte, n)

	/* NOTE(anton2920): see encoding/base64/base64.go:294. */
	token := make([]byte, (n+2)/3*4)

	/* TODO(anton2920): think about adding delay before next attempt. */
	for {
		/* NOTE(anton2920): DEV0-28. Usage of current time as a part of the token must prevent generation of identical tokens for different IDs. */
		t := time.Unix()
		*(*int64)(unsafe.Pointer(&buffer[0])) = t

		n, err := rand.Read(buffer[unsafe.Sizeof(t):])
		if (err != nil) || (n != len(buffer[unsafe.Sizeof(t):])) {
			log.Warnf("Failed to read random bytes: %v", err)
			continue
		}

		base64.StdEncoding.Encode(token, buffer)

		/* Making sure that it's unique. */
		SessionsLock.RLock()
		_, ok := Sessions[bytes.AsString(buffer)]
		SessionsLock.RUnlock()
		if !ok {
			return string(token)
		}
	}
}

func LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open sessions file %q: %v", filename, err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(&Sessions); err != nil {
		return fmt.Errorf("failed to decode sessions from file: %v", err)
	}

	return nil
}

func StoreToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create sessions file %q: %v", filename, err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	SessionsLock.Lock()
	defer SessionsLock.Unlock()
	if err := enc.Encode(Sessions); err != nil {
		return fmt.Errorf("failed to encode sessions to file: %v", err)
	}

	return nil
}
