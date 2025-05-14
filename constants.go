package pastebin

const (
	BaseUrl      = "https://pastebin.com"
	LoginUrl     = "https://pastebin.com/api/api_login.php"
	PostUrl      = "https://pastebin.com/api/api_post.php"
	RawUrl       = "https://pastebin.com/api/api_raw.php"
	RawPublicUrl = "https://pastebin.com/raw"
)

const (
	Never      Expiration = "N"
	TenMinutes Expiration = "10M"
	OneHour    Expiration = "1H"
	OneDay     Expiration = "1D"
	OneWeek    Expiration = "1W"
	TwoWeeks   Expiration = "2W"
	OneMonth   Expiration = "1M"
	SixMonths  Expiration = "6M"
	OneYear    Expiration = "1Y"
)

const (
	Public Visibility = iota
	Unlisted
	Private
)

func (v Visibility) String() string {
	switch v {
	case Public:
		return "Public"
	case Unlisted:
		return "Unlisted"
	case Private:
		return "Private"
	default:
		return "Unknown"
	}
}

const (
	NormalUser AccountType = iota
	ProUser
)

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

const (
	apiUserName        = "api_user_name"
	apiUserKey         = "api_user_key"
	apiUserPassword    = "api_user_password"
	apiDevKey          = "api_dev_key"
	apiOption          = "api_option"
	apiPasteCode       = "api_paste_code"
	apiPasteName       = "api_paste_name"
	apiPasteKey        = "api_paste_key"
	apiPasteFormat     = "api_paste_format"
	apiPasteExpireDate = "api_paste_expire_date"
	apiPastePrivate    = "api_paste_private"
	apiFolderKey       = "api_folder_key"
	apiResultsLimit    = "api_results_limit"
)
