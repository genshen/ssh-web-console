package utils

import (
	"time"
	"sync"
)

var SessionStorage SessionManager

var mutex = new(sync.RWMutex)

func initSessionUtils() {
	SessionStorage.new()
}

// use jwt string as session key,
// store user information(username and password) in Session.
type SessionManager struct {
	sessions map[string]Session
}

type Session struct {
	expire int64
	Value  interface{}
}

func (s *Session) isExpired(timeNow int64) bool {
	if s.expire < timeNow {
		return true
	}
	return false
}

func (s *SessionManager) new() {
	s.sessions = make(map[string]Session)
}

/**
* add a new session to session manager.
* @params:token: token string
* expire: unix time for expire
* password: ssh user password
 */
func (s *SessionManager) Put(key string, expire int64, value interface{}) {
	s.gc()
	mutex.Lock()
	s.sessions[key] = Session{expire: expire, Value: value}
	mutex.Unlock()
}

func (s *SessionManager) Get(key string) (sessionData Session, exist bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	session, ok := s.sessions[key]
	return session, ok
}

func (s *SessionManager) Delete(key string) {
	mutex.Lock()
	if _, ok := s.sessions[key]; ok {
		delete(s.sessions, key)
	}
	mutex.Unlock()
}

func (s *SessionManager) gc() {
	timeNow := time.Now().Unix()
	mutex.Lock()
	for key, session := range s.sessions {
		if session.isExpired(timeNow) {
			delete(s.sessions, key)
		}
	}
	mutex.Unlock()
}
