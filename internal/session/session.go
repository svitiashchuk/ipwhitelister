package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type SessionManager struct {
	sessions map[string]string
	mutex    sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]string),
		mutex:    sync.Mutex{},
	}
}

func (sm *SessionManager) StoreSessionData(sessionToken, value string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.sessions[sessionToken] = value
}

func (sm *SessionManager) RetrieveSessionData(sessionToken string) (string, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	value, ok := sm.sessions[sessionToken]

	return value, ok
}

func (sm *SessionManager) DeleteSessionData(sessionToken string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.sessions, sessionToken)
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
