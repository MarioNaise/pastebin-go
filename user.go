package pastebin

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
