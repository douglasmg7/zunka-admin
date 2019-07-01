package main

import (
	// _ "github.com/mattn/go-sqlite3"
	// "github.com/satori/go.uuid"

	"database/sql"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Measure execution time to retrive the session.
// var timeToGetSession time.Time

// SessionData cache data from each user.
type SessionData struct {
	UserID     int
	UserName   string
	Permission uint64
	Outdated   bool // If true, session must be retrived from db.
}

// CheckPermission return if session has the permission.
func (s *SessionData) CheckPermission(p string) bool {
	// Admin.
	if s.Permission&1 == 1 {
		return true
	}
	// Permissions.
	switch p {
	case "editStudent":
		return s.Permission&2 == 2
	case "editPrice":
		return s.Permission&4 == 4
	default:
		return false
	}
}

// SetPermission grand permission.
func (s *SessionData) SetPermission(p string) {
	switch p {
	case "admin":
		s.Permission = s.Permission | 1
	case "editStudent":
		s.Permission = s.Permission | 2
	case "editPrice":
		s.Permission = s.Permission | 4
	}
}

// UnsetPermission revoke permission.
func (s *SessionData) UnsetPermission(p string) {
	switch p {
	case "king":
		s.Permission = s.Permission ^ 1
	case "editStudent":
		s.Permission = s.Permission ^ 2
	case "editPrice":
		s.Permission = s.Permission ^ 4
	}
}

// PasswordIsCorrect return if password is correct.
func (s *SessionData) PasswordIsCorrect(password string) bool {
	// Get user by email.
	var cryptedPassword []byte
	err = db.QueryRow("SELECT password FROM user WHERE id = ?", s.UserID).Scan(&cryptedPassword)
	// No registred user.
	if err == sql.ErrNoRows {
		return false
	}
	// Internal error.
	if err != nil {
		log.Fatal(err)
	}
	// Compare password.
	err = bcrypt.CompareHashAndPassword(cryptedPassword, []byte(password))
	// Incorrect password.
	if err != nil {
		return false
	}
	// Correct password.
	return true
}

// Sessions contains sessions from each user and userId from each uuid sesscion.
type Sessions struct {
	// UserId from uuidSession.
	mapUserID map[string]int
	// Session from userId.
	mapSessionData map[int]*SessionData
}

// CreateSession create a session and writing a cookie on client and keep a reletion of cookie -> user id.
func (s *Sessions) CreateSession(w http.ResponseWriter, userID int) error {
	// create cookie
	sUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	sUUIDString := sUUID.String()
	// Save cookie.
	http.SetCookie(w, &http.Cookie{
		Name:  "sessionUUID",
		Value: sUUIDString,
		Path:  "/",
		// Secure: true, // to use only in https
		// HttpOnly: true, // Can't be used into js client
	})
	// Save session UUID on db.
	stmt, err := db.Prepare(`INSERT INTO sessionUUID(uuid, user_id, createdAt) VALUES( ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(sUUIDString, userID, time.Now())
	if err != nil {
		return err
	}
	// Save on cache.
	s.mapUserID[sUUIDString] = userID
	return nil
}

// RemoveSession remove session from client browser.
func (s *Sessions) RemoveSession(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("sessionUUID")
	// No cookie.
	if err == http.ErrNoCookie {
		// log.Println("No cookie")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		// Some error.
	} else if err != nil {
		log.Fatal(err)
		// Remove cookie.
	} else {
		c.MaxAge = -1
		c.Path = "/"
		// log.Println("changed cookie:", c)
		http.SetCookie(w, c)
		http.Redirect(w, req, "/auth/signin", http.StatusSeeOther)
		// Delete userId session.
		delete(s.mapUserID, c.Value)
	}
}

// GetSession return session data from UUID.
func (s *Sessions) GetSession(req *http.Request) (*SessionData, error) {
	// timeToGetSession = time.Now()
	userID, err := s.getUserIdfromSessionUUID(req)
	// Some error.
	if err != nil {
		return nil, err
		// No user id.
	} else if userID == 0 {
		return nil, nil
		// Found user.
	} else {
		session, err := s.getSessionFromUserID(userID)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println("Time to get session:", time.Since(timeToGetSession))
		return session, err
		// return sessionDataFromUserId(userID)
	}
}

// Return user id from session uuid.
// Try the cache first.
func (s *Sessions) getUserIdfromSessionUUID(req *http.Request) (int, error) {
	cookie, err := req.Cookie("sessionUUID")
	// log.Println("Cookie:", cookie.Value)
	// log.Println("Cookie-err:", err)
	// No cookie.
	if err == http.ErrNoCookie {
		return 0, nil
		// some error
	} else if err != nil {
		return 0, err
	}
	// Have a cookie.
	if cookie != nil {
		sessionUUID := cookie.Value
		userID := s.mapUserID[sessionUUID]
		// Found on cache.
		if userID != 0 {
			// log.Println("userId from cache", userId)
			return userID, nil
		}
		// Get from db.
		err = db.QueryRow("select user_id from sessionUUID where uuid = ?", sessionUUID).Scan(&userID)
		if err == sql.ErrNoRows {
			return 0, nil
		}
		if err != nil {
			// No user id for the sessionUUID.
			return 0, err
		}
		// Found the user id.
		if userID != 0 {
			// log.Println("userId from db", userId)
			s.mapUserID[sessionUUID] = userID
			return userID, nil
		}
	}
	// No cookie
	return 0, nil
}

// Get session from cache.
// If not cached, cache it.
func (s *Sessions) getSessionFromUserID(userID int) (session *SessionData, err error) {
	// From the cache.
	session = s.mapSessionData[userID]
	// No cached nor outdated.
	if session != nil && !session.Outdated {
		return session, nil
	}
	// Cache from db.
	// log.Println("Getting session data from cache:")
	return s.cacheSession(userID)
}

// Cache session data and return it.
func (s *Sessions) cacheSession(userID int) (session *SessionData, err error) {
	session = &SessionData{}
	err = db.QueryRow("select id, name, permission from user where id = ?", userID).Scan(&session.UserID, &session.UserName, &session.Permission)
	if err != nil {
		return nil, err
	}
	// Cache it.
	s.mapSessionData[userID] = session
	return session, nil
}

// CleanSessions clean the cache session.
func (s *Sessions) CleanSessions() {
	s.mapSessionData = map[int]*SessionData{}
	log.Println("Sessions cache cleaned")
}
