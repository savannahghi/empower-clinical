package domain

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
)

// This file mirrors the models that are defined in myCareHub.

// User  holds details that both the client and staff have in common
//
// Client and Staff cannot exist without being a user
type User struct {
	ID          *string          `json:"userID"`
	Username    string           `json:"userName"`
	UserType    string           `json:"userType"`
	Name        string           `json:"name"`
	Gender      enumutils.Gender `json:"gender"`
	Active      bool             `json:"active"`
	Contacts    *Contact         `json:"primaryContact"`
	Flavour     feedlib.Flavour  `json:"flavour"`
	Avatar      string           `json:"avatar"`
	DateOfBirth *time.Time       `json:"dateOfBirth"`
}

// Contact hold contact information/details for users
type Contact struct {
	ID           *string `json:"id"`
	ContactType  string  `json:"contactType"`
	ContactValue string  `json:"contactValue"`
	Active       bool    `json:"active"`
	OptedIn      bool    `json:"optedIn"`
}
