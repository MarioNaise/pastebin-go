package pastebin

import "fmt"

// AccountType represents a user account type on Pastebin.
// normal User = 0, pro User = 1
//
// See https://pastebin.com/doc_api#12
type AccountType int

// String returns the string representation of an AccountType.
func (acc AccountType) String() string {
	switch AccountType(acc) {
	case NormalUser:
		return "NormalUser"
	case ProUser:
		return "ProUser"
	default:
		return "UnknownUserType"
	}
}

// User contains information about the logged in Pastebin user.
type User struct {
	UserName    string      `xml:"user_name"`
	Expiration  Expiration  `xml:"user_expiration"`
	Visibility  Visibility  `xml:"user_private"`
	Avatar      string      `xml:"user_avatar_url"`
	Website     string      `xml:"user_website"`
	Email       string      `xml:"user_email"`
	Location    string      `xml:"user_location"`
	AccountType AccountType `xml:"user_account_type"`
}

// String returns a formatted string of the user data.
func (u User) String() string {
	return fmt.Sprintf("UserName: %s, Expiration: %s, Visibility: %s, Avatar: %s, Website: %s, Email: %s, Location: %s, AccountType: %s",
		u.UserName, u.Expiration, u.Visibility, u.Avatar, u.Website, u.Email, u.Location, u.AccountType)
}
