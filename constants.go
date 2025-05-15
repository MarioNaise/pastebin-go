package pastebin

// Base URLs and API endpoints for Pastebin.
const (
	BaseUrl      = "https://pastebin.com"
	LoginUrl     = "https://pastebin.com/api/api_login.php"
	PostUrl      = "https://pastebin.com/api/api_post.php"
	RawUrl       = "https://pastebin.com/api/api_raw.php"
	RawPublicUrl = "https://pastebin.com/raw"
)

// Predefined expiration times for pastes.
//
// See https://pastebin.com/doc_api#6
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

// Public = 0, Unlisted = 1, Private = 2.
//
// See https://pastebin.com/doc_api#7
const (
	Public Visibility = iota
	Unlisted
	Private
)

// NormalUser is a free Pastebin account.
// ProUser is a paid Pastebin account.
//
// See https://pastebin.com/doc_api#12
const (
	NormalUser AccountType = iota
	ProUser
)

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
