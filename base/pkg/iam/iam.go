package iam

import (
	"context"
	"errors"
)

type Scope string

var (
	ErrTokenExpired            = errors.New("access token is expired")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrGroupNotFound           = errors.New("group not found")
	//ErrGroupAlreadyExists      = errors.New("group already exists")
	NullPermissions = []string{} // const value signifying no particular permissions required
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Enabled   bool   `json:"enabled"`
}

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IAM interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Logout(ctx context.Context, token string) error
	Authorize(ctx context.Context, token string, permissions []string) (userID string, err error)
	SetPassword(username string, password string, temporary bool) error

	CreateUser(user *User, password string, temporary bool) (id string, err error)
	DeleteUser(username string) error
	GetUser(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(context.Context) ([]*User, error)
	SearchUsers(ctx context.Context, prefix string) ([]*User, error)

	CreateGroup(name string) (*Group, error)
	DeleteGroup(groupID string) error
	DeleteGroupByName(name string) error
	AddUserToGroup(userID string, groupID string) error
	RemoveUserFromGroup(userID string, groupID string) error
	ListUserGroups(ctx context.Context, userID string) ([]*Group, error)
	ListGroupMembers(ctx context.Context, groupID string) ([]*User, error)
	IsUserInGroup(ctx context.Context, userID string, groupID string) (bool, error)
	ResolveGroup(ctx context.Context, groupName string) (groupID string, err error)
}
