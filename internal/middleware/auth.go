package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
	websocket "forum/internal/ws"
)

type visitor struct {
	requests []time.Time
	lastSeen time.Time
	mu       sync.Mutex
}

type RateLimiter struct {
	visitors    map[string]*visitor
	mu          sync.RWMutex
	maxRequests int
	window      time.Duration
	cleanup     time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:    make(map[string]*visitor),
		maxRequests: maxRequests,
		window:      window,
		cleanup:     time.Hour,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanup)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			requests: make([]time.Time, 0, rl.maxRequests),
		}
		rl.visitors[ip] = v
	}
	rl.mu.Unlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	v.lastSeen = now

	windowStart := now.Add(-rl.window)
	valid := 0
	for _, t := range v.requests {
		if t.After(windowStart) {
			v.requests[valid] = t
			valid++
		}
	}
	v.requests = v.requests[:valid]

	if len(v.requests) < rl.maxRequests {
		v.requests = append(v.requests, now)
		return true
	}
	return false
}

type customHandler func(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int, hub *websocket.Hub)

func AuthMiddleware(db *sql.DB, next customHandler, maxRequests int, window time.Duration, hub *websocket.Hub, loged bool) http.Handler {
	rateLimiter := NewRateLimiter(maxRequests, window)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get IP address
		ip := r.RemoteAddr
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ip = forwardedFor
		}

		// Check rate limit
		if !rateLimiter.isAllowed(ip) {
			utils.RespondWithJSON(w, http.StatusTooManyRequests, utils.ErrorResponse{
				Error: "Rate limit exceeded. Please try again later.",
			})
			return
		}
		// validation session
		if loged {
			userId, err := ValidUser(r, db)
			if err != nil {
				if err == http.ErrNoCookie {
					// no session in database new user
					utils.RespondWithJSON(w, http.StatusUnauthorized, utils.ErrorResponse{Error: "Unauthorized"})
					return
				} else if err == sql.ErrNoRows {
					// user with expirated date we clean the last session cookie
					http.SetCookie(w, &http.Cookie{
						Name:    "session_token",
						Path:    "/",
						Value:   "",
						Expires: time.Unix(0, 0),
					})
					utils.RespondWithJSON(w, http.StatusUnauthorized, utils.ErrorResponse{Error: "Unauthorized"})
					return
				} else {
					utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
					return
				}
			}
			next(w, r, db, userId, hub)
		} else {
			next(w, r, db, 0, nil)
		}
	})
}

func IsUserRegistered(db *sql.DB, userData *utils.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?);`
	err := db.QueryRow(query, userData.Email, userData.UserName).Scan(&exists)
	return exists, err
}

func RegisterUser(db *sql.DB, userData *utils.User) error {
	insertQuery := `INSERT INTO users (username, Age, Gender, First_Name, Last_Name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?);`
	result, err := db.Exec(insertQuery, userData.UserName, userData.Age, userData.Gender, userData.FirstName, userData.LastName, userData.Email, userData.Password)
	if err != nil {
		return err
	}
	userData.UserId, err = result.LastInsertId()
	return err
}

func GetActiveSession(db *sql.DB, userData *utils.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id = ?  AND expires_at > ?);`
	err := db.QueryRow(query, userData.UserId, userData.Expiration).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func DeleteSession(db *sql.DB, userData *utils.User) error {
	query := `DELETE FROM sessions WHERE user_id =  ?;`
	_, err := db.Exec(query, userData.UserId)
	return err
}

func ValidCredential(db *sql.DB, userData *utils.User) error {
	query := `SELECT id, password FROM users WHERE (username = ? OR email= ?);`
	err := db.QueryRow(query, userData.UserName, userData.Email).Scan(&userData.UserId, &userData.Password)
	if err != nil {
		fmt.Println("test test")
		return err
	}
	return err
}

func ValidUser(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}
	userid, err := database.Get_session(cookie.Value, db)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func RemoveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cookie.Value)
	if err != nil {
		return err
	}
	return nil
}
