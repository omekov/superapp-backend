package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// User ...
type User struct {
	ID           uuid.UUID `db:"id"`
	UserName     string    `db:"username"`
	Password     string    `db:"password"`
	Email        string    `db:"email"`
	State        string    `db:"state"`
	PinCode      string    `db:"pin_code"`
	Address      Address
	Profile      Profile `db:"profile"`
	Phone        Phone
	IdentityCard IdentityCard
	Group        Group
	Status       Status
}

// TestGetUser ...
func TestGetUser(t *testing.T) User {
	t.Helper()
	id, _ := uuid.Parse("280fe6bc-7fea-48c2-bf25-953408d0879f")
	return User{
		ID:       id,
		UserName: "azamat",
		Password: "password",
		State:    "enabled",
		Email:    "azamat@example.com",
	}
}

// TestGetUser2 ...
func TestGetUser2(t *testing.T) User {
	t.Helper()
	id, _ := uuid.Parse("50b0e21d-1a10-465b-8515-df42b281f251")
	return User{
		ID:       id,
		UserName: "azamat2",
		Password: "password",
		State:    "enabled",
		Email:    "azamat2@example.com",
	}
}

// TestGetUser3 ...
func TestGetUser3(t *testing.T) User {
	t.Helper()
	id, _ := uuid.Parse("fa222a11-6eb9-4540-b265-df99a7f9314b")
	return User{
		ID:       id,
		UserName: "azamat3",
		Password: "password",
		State:    "enabled",
		Email:    "azamat3@example.com",
	}
}

// TestGetUser4 ...
func TestGetUser4(t *testing.T) User {
	t.Helper()
	id, _ := uuid.Parse("42cf6af9-ce43-4e44-8f29-c8890eb1f9ed")
	return User{
		ID:       id,
		UserName: "azamat4",
		Password: "password",
		State:    "enabled",
		Email:    "azamat4@example.com",
	}
}

// TestGetUser5 ...
func TestGetUser5(t *testing.T) User {
	t.Helper()
	id, _ := uuid.Parse("41b4b7a6-bdf6-47a0-9b5d-8c62693ddaf0")
	return User{
		ID:       id,
		UserName: "azamat5",
		Password: "password",
		State:    "enabled",
		Email:    "azamat5@example.com",
	}
}

// TestGetUser6 ...
func TestGetUser6(t *testing.T) User {
	t.Helper()
	id, _ := uuid.Parse("4cc54cc6-92fb-4002-bf3c-331b8762eac0")
	return User{
		ID:       id,
		UserName: "azamat6",
		Password: "password",
		State:    "enabled",
		Email:    "azamat6@example.com",
		PinCode:  "654321",
	}
}

type Address struct {
	Country  string
	Region   string
	City     string
	District string
	Street   string
	House    string
	Room     string
	Lat      string
	Lng      string
}

type Profile struct {
	FirstName  string
	LastName   string
	MiddleName string
	BirthDate  time.Time
	Sex        uint
}

type Phone struct {
	Number uint
}

type IdentityCard struct {
	Issuer         string
	Number         string
	IssueDate      time.Time
	ExpirationDate time.Time
}

type Status struct {
	ID      uuid.UUID
	NameEng string
	NameQaz string
	NameRus string
	Code    string
}

type Role struct {
	ID         uuid.UUID
	NameEng    string
	NameQaz    string
	NameRus    string
	Code       string
	Permission []Permission
}

type Group struct {
	ID      uuid.UUID
	NameEng string
	NameQaz string
	NameRus string
	Code    string
	Role    []Role
}

type Permission struct {
	ID      uuid.UUID
	NameEng string
	NameQaz string
	NameRus string
	Code    string
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
}
