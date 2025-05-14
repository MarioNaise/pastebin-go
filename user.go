package pastebin

import "fmt"

type AccountType int

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

func (u User) String() string {
	return fmt.Sprintf("UserName: %s, Expiration: %s, Visibility: %s, Avatar: %s, Website: %s, Email: %s, Location: %s, AccountType: %s",
		u.UserName, u.Expiration, u.Visibility, u.Avatar, u.Website, u.Email, u.Location, u.AccountType)
}
