package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

var (
	minUsernameLength = 4
	minPasswordLength = 8
	reservedUsernames = []string{"admin", "sysadmin", "administrator"}
	reservedEmails    = []string{"d@d.com"}
)

func (s *Server) usernameTaken(
	ctx context.Context,
	name string,
) (bool, error) {
	if base.Contains(reservedUsernames, name) {
		return true, nil
	}
	if _, err := s.iam.GetUser(
		ctx,
		name,
	); err == iam.ErrUserNotFound {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "iam")
	} else {
		return true, nil
	}
}

func (s *Server) emailTaken(
	ctx context.Context,
	email string,
) (bool, error) {
	if base.Contains(reservedEmails, email) {
		return true, nil
	}
	if _, err := s.iam.GetUserByEmail(
		ctx,
		email,
	); err == iam.ErrUserNotFound {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "iam")
	} else {
		return true, nil
	}
}

func (s *Server) handleUserExists() http.HandlerFunc {
	return s.handler(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) (err error) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return nil
			}
			var result struct {
				Exists bool `json:"exists"`
			}
			if username := r.URL.Query().Get("u"); username != "" {
				if r.URL.Query().Has("e") {
					writeError(
						w,
						http.StatusBadRequest,
						errors.New("username and email are mutually exclusive"),
					)
					return nil
				}
				if result.Exists, err = s.usernameTaken(
					r.Context(),
					username,
				); err != nil {
					return err
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(&result); err != nil {
					return errors.Wrap(err, "json")
				}
				s.log.Debug("check username taken",
					zap.String("username", username),
					zap.Bool("exists", result.Exists))
				return nil
			}
			if email := r.URL.Query().Get("e"); email != "" {
				if result.Exists, err = s.emailTaken(
					r.Context(),
					email,
				); err != nil {
					return err
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(&result); err != nil {
					return errors.Wrap(err, "json")
				}
				s.log.Debug("check email taken",
					zap.String("email", email),
					zap.Bool("exists", result.Exists))
				return nil
			}
			writeError(
				w,
				http.StatusBadRequest,
				errors.New("username or email is required"),
			)
			return nil
		})
}

func (s *Server) handleSetPassword() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			// TODO: get requesting user details
			var req struct {
				NewPassword string `json:"newPassword"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "json decode")
			}
			return nil
		})
}

func (s *Server) handleRegister() http.HandlerFunc {
	return s.handler(
		http.MethodPost,
		func(w http.ResponseWriter, r *http.Request) (err error) {
			var req struct {
				Username  string `json:"username"`
				Email     string `json:"email"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Password  string `json:"password"`
			}
			if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "json decode")
			}
			req.Username = strings.TrimSpace(req.Username)
			req.Email = strings.TrimSpace(req.Email)
			req.FirstName = strings.TrimSpace(req.FirstName)
			req.LastName = strings.TrimSpace(req.LastName)
			s.log.Debug("registering user",
				zap.String("username", req.Username),
				zap.String("email", req.Email))
			if !isValidUsername(req.Username) {
				return errors.New("invalid username")
			}
			if !isValidPassword(req.Password) {
				return errors.New("password does not meet requirements")
			}
			if !isValidEmail(req.Email) {
				return errors.New("invalid email")
			}
			if req.FirstName == "" {
				return errors.New("missing first name")
			}
			if req.LastName == "" {
				return errors.New("missing last name")
			}
			var res struct {
				*iam.User   `json:""`
				AccessToken string `json:"accessToken"`
			}
			user := &iam.User{
				Username:  req.Username,
				Email:     req.Email,
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Enabled:   true,
			}
			res.User = user
			if res.ID, err = s.iam.CreateUser(
				user,
				req.Password,
				false,
			); err != nil {
				return errors.Wrap(err, "iam.CreateUser")
			}
			if res.AccessToken, err = s.iam.Login(
				r.Context(),
				req.Username,
				req.Password,
			); err != nil {
				return errors.Wrap(err, "iam.Login")
			}
			if resUserID, err := s.iam.Authorize(
				r.Context(),
				res.AccessToken,
				nil,
			); err != nil {
				return errors.Wrap(err, "iam.Authorize")
			} else if resUserID != user.ID {
				return fmt.Errorf("failed id sanity check (got '%s', expected '%s')",
					resUserID,
					user.ID,
				)
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(&res); err != nil {
				return errors.Wrap(err, "json encode")
			}
			s.log.Debug("registered user",
				zap.String("user.Username", res.Username),
				zap.String("user.Email", res.Email),
				zap.String("user.FirstName", res.FirstName),
				zap.String("user.LastName", res.LastName),
				zap.Bool("user.Enabled", res.Enabled))
			return nil
		})
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidUsername(username string) bool {
	if len(username) < minUsernameLength {
		return false
	} else if base.Contains(reservedUsernames, username) {
		return false
	}
	return true
}

func isValidPassword(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}
	b := []byte(password)
	for _, pattern := range []string{
		"[a-z]",
		"[A-Z]",
		"[0-9]",
		"[$&+,:;=?@#|'<>.^*()%!-]",
	} {
		if match, _ := regexp.Match(pattern, b); !match {
			return false
		}
	}
	return true
}

func (s *Server) handleLogin() http.HandlerFunc {
	return s.handler(
		http.MethodPost,
		func(w http.ResponseWriter, r *http.Request) error {
			var req struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "json decode")
			}
			s.log.Debug("logging in",
				zap.String("username", req.Username))
			var response struct {
				*iam.User   `json:""`
				AccessToken string `json:"accessToken"`
			}
			var err error
			response.AccessToken, err = s.iam.Login(
				r.Context(),
				req.Username,
				req.Password,
			)
			if err == iam.ErrInvalidCredentials {
				w.WriteHeader(http.StatusUnauthorized)
				s.log.Warn("failed login attempt",
					zap.String("username", req.Username))
				return nil
			} else if err != nil {
				return errors.Wrap(err, "iam.Login")
			}
			response.User, err = s.iam.GetUser(r.Context(), req.Username)
			if err != nil {
				return errors.Wrap(err, "iam.GetUser")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(&response); err != nil {
				return errors.Wrap(err, "json encode")
			}
			return nil
		})
}

func (s *Server) handleSignOut() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			return s.iam.Logout(
				r.Context(),
				r.Header.Get("AccessToken"),
			)
		})
}

func (s *Server) handleUserSearch() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			prefix := r.URL.Query().Get("p")
			if len(prefix) < minUsernameLength {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			users, err := s.iam.SearchUsers(r.Context(), prefix)
			if err != nil {
				return errors.Wrap(err, "iam.SearchUsers")
			}
			response := make([]*searchUser, len(users))
			for i, user := range users {
				response[i] = &searchUser{
					ID:       user.ID,
					Username: user.Username,
				}
			}
			w.Header().Set("Content-Type", "application/json")
			s.log.Debug("searched for user",
				zap.String("prefix", prefix),
				zap.Int("numResults", len(response)))
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(&response); err != nil {
				return errors.Wrap(err, "json encode")
			}
			return nil
		})
}

type searchUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
