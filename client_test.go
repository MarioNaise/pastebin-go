package pastebin

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TODO: clean up tests
// TODO: create test package
type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func getMockHttpClient(resp *http.Response) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return resp
		}),
	}
}

type roundTripFuncErr func(req *http.Request) *http.Response

func (f roundTripFuncErr) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("mock error")
}

func getMockHttpClientErr() *http.Client {
	return &http.Client{
		Transport: roundTripFuncErr(func(req *http.Request) *http.Response {
			return nil
		}),
	}
}

func TestNewClient(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if client.apiUserName != "testUserName" {
		t.Error("Expected apiUserName to be testUserName, got", client.apiUserName)
	}
	if client.apiUserPassword != "testPassword" {
		t.Error("Expected apiUserPassword to be testPassword, got", client.apiUserPassword)
	}
	if client.apiDevKey != "testApiKey" {
		t.Error("Expected apiDevKey to be testApiKey, got", client.apiUserPassword)
	}
}

func TestCreatePaste(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("https://pastebin.com/pasteKey")),
	})
	key, err := client.CreatePaste(&CreatePasteRequest{})
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if key != "pasteKey" {
		t.Errorf("Expected pasteKey, got %s", key)
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	_, err = client.CreatePaste(&CreatePasteRequest{})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	values := url.Values{}
	httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Error("Expected no error reading body, got", err)
			}
			encodedValue := string(body)
			values, err = url.ParseQuery(encodedValue)
			if err != nil {
				t.Error("Expected no error parsing query, got", err)
			}
			return &http.Response{Body: io.NopCloser(strings.NewReader("https://pastebin.com/testKey"))}
		}),
	}
	res, err := client.CreatePaste(&CreatePasteRequest{Content: "testContent", Folder: "go", Visibility: Unlisted, CreatePasteAsUser: true})
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if values.Get(apiPasteCode) != "testContent" {
		t.Errorf("Expected api_paste_code to be testContent, got %s", values.Get(apiPasteCode))
	}
	if values.Get(apiOption) != "paste" {
		t.Errorf("Expected api_option to be paste, got %s", values.Get("api_option"))
	}
	if values.Get(apiDevKey) != "testApiKey" {
		t.Errorf("Expected api_dev_key to be testApiKey, got %s", values.Get("api_dev_key"))
	}
	if values.Get(apiUserKey) != "testUserKey" {
		t.Errorf("Expected api_user_key to be testUserKey, got %s", values.Get("api_user_key"))
	}
	if values.Get(apiFolderKey) != "go" {
		t.Errorf("Expected api_folder_key to be go, got %s", values.Get("api_folder_key"))
	}
	if values.Get(apiPastePrivate) != "1" {
		t.Errorf("Expected api_visibility to be 1, got %s", values.Get("api_visibility"))
	}
	if res != "testKey" {
		t.Errorf("Expected testKey, got %s", res)
	}
}

func TestGetUserPastes(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader(`<paste>
			<paste_key>0b42rwhf</paste_key>
			<paste_date>1297953260</paste_date>
			<paste_title>javascript test</paste_title>
			<paste_size>15</paste_size>
			<paste_expire_date>1297956860</paste_expire_date>
			<paste_private>0</paste_private>
			<paste_format_long>JavaScript</paste_format_long>
			<paste_format_short>javascript</paste_format_short>
			<paste_url>https://pastebin.com/0b42rwhf</paste_url>
			<paste_hits>15</paste_hits>
		</paste>
		<paste>
			<paste_key>0C343n0d</paste_key>
			<paste_date>1297694343</paste_date>
			<paste_title>Welcome To Pastebin V3</paste_title>
			<paste_size>490</paste_size>
			<paste_expire_date>0</paste_expire_date>
			<paste_private>0</paste_private>
			<paste_format_long>None</paste_format_long>
			<paste_format_short>text</paste_format_short>
			<paste_url>https://pastebin.com/0C343n0d</paste_url>
			<paste_hits>65</paste_hits>
		</paste>`)),
	})

	pastes, err := client.GetUserPastes()
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if len(pastes) != 2 {
		t.Errorf("Expected 2 pastes, got %d", len(pastes))
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	_, err = client.GetUserPastes()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("<test>")),
	})
	_, err = client.GetUserPastes()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDeletePaste(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Paste Removed")),
	})
	if err = client.DeletePaste("testPasteKey"); err != nil {
		t.Error("Expected no error, got", err)
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	if err = client.DeletePaste("testPasteKey"); err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("unexpected response")),
	})
	if err = client.DeletePaste("testPasteKey"); err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetRawUserPasteContent(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("test user paste")),
	})
	res, err := client.GetRawUserPasteContent("testPasteKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if res != "test user paste" {
		t.Errorf("Expected 'test user paste', got '%s'", res)
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	_, err = client.GetRawUserPasteContent("testPasteKey")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetRawPublicUserPasteContent(t *testing.T) {
	client, err := NewClient("", "", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("test public paste")),
	})
	res, err := client.GetRawPublicPasteContent("testPasteKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if res != "test public paste" {
		t.Error("Expected 'test public paste', got", res)
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	_, err = client.GetRawPublicPasteContent("testPasteKey")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = getMockHttpClientErr()
	_, err = client.GetRawPublicPasteContent("testPasteKey")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				Body: io.NopCloser(strings.NewReader("Bad API request, invalid api_user_key")),
			}
		}),
	}
	_, err = client.GetRawPublicPasteContent("testPasteKey")
	if err == nil {
		t.Error("Expected no error, got", err)
	}

	httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			resp := "Bad API request"
			if req.URL.String() == RawPublicUrl+"/testPasteKey" {
				resp = "Bad API request, invalid api_user_key"
				if client.apiUserKey == "testUserKey" {
					resp = "testResponse"
				}
				return &http.Response{
					Body: io.NopCloser(strings.NewReader(resp)),
				}
			} else if req.URL.String() == LoginUrl {
				return &http.Response{
					Body: io.NopCloser(strings.NewReader("testUserKey")),
				}
			}
			return &http.Response{
				Body: io.NopCloser(strings.NewReader(resp)),
			}
		}),
	}
	r, err := client.GetRawPublicPasteContent("testPasteKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if r != "testResponse" {
		t.Errorf("Expected testResponse, got %s", r)
	}
}

func TestGetUserDetails(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader(`<user>
        <user_name>wiz_kitty</user_name>
        <user_format_short>text</user_format_short>
        <user_expiration>N</user_expiration>
        <user_avatar_url>https://pastebin.com/cache/a/1.jpg</user_avatar_url>
        <user_private>1</user_private>
        <user_website>https://myawesomesite.com</user_website>
        <user_email>oh@dear.com</user_email>
        <user_location>New York</user_location>
        <user_account_type>1</user_account_type>
</user>`)),
	})
	user, err := client.GetUserDetails()
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if user.UserName != "wiz_kitty" {
		t.Error("Expected userName to be wiz_kitty, got", user.UserName)
	}
	if user.Visibility != Unlisted {
		t.Error("Expected visibility to be Unlisted, got", user.Visibility)
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	_, err = client.GetUserDetails()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("unexpected response")),
	})
	_, err = client.GetUserDetails()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRefreshUserKey(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testUserKey")),
	})
	client, err := NewClient("testUserName", "testPassword", "testApiKey")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if client.apiUserKey != "testUserKey" {
		t.Error("Expected apiUserKey to be testUserKey, got", client.apiUserKey)
	}

	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("Bad API request")),
	})
	if err = client.refreshUserKey(); err == nil {
		t.Error("Expected error, got nil")
	}
	if client.apiUserKey != "testUserKey" {
		t.Error("Expected apiUserKey to be testUserKey, got", client.apiUserKey)
	}
}

func TestPost(t *testing.T) {
	httpClient = getMockHttpClient(&http.Response{
		Body: io.NopCloser(strings.NewReader("testResult")),
	})
	client, _ := NewClient("", "", "")
	r, err := client.post("https://example.com", nil, false)
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if r != "testResult" {
		t.Errorf("Expected testResult, got %s", r)
	}

	_, err = client.post("://example.com", nil, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = getMockHttpClientErr()
	_, err = client.post("https://example.com", nil, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				Body: io.NopCloser(strings.NewReader("Bad API request, invalid api_user_key")),
			}
		}),
	}
	_, err = client.post("https://example.com", nil, true)
	if err == nil {
		t.Error("Expected no error, got", err)
	}

	httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			if req.URL.String() == "https://example.com" {
				resp := "Bad API request, invalid api_user_key"
				if client.apiUserKey == "testUserKey" {
					resp = "testResponse"
				}
				return &http.Response{
					Body: io.NopCloser(strings.NewReader(resp)),
				}
			}
			return &http.Response{
				Body: io.NopCloser(strings.NewReader("testUserKey")),
			}
		}),
	}
	r, err = client.post("https://example.com", nil, true)
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if r != "testResponse" {
		t.Errorf("Expected testResponse, got %s", r)
	}
}

func TestNewRequest(t *testing.T) {
	result, err := newRequest("https://example.com/api", url.Values{
		"api_test_key": {"test"},
	})
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if result == nil {
		t.Error("Expected non-nil request, got nil")
	}
	bodyContent, err := io.ReadAll(result.Body)
	if err != nil {
		t.Error("Expected no error reading body, got", err)
	}
	if string(bodyContent) != "api_test_key=test" {
		t.Error("Expected body content 'api_test_key=test', got", string(bodyContent))
	}
	if result.Method != "POST" {
		t.Error("Expected POST method, got", result.Method)
	}
	if result.URL.String() != "https://example.com/api" {
		t.Error("Expected URL https://example.com/api, got", result.URL.String())
	}
	if result.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Error("Expected Content-Type application/x-www-form-urlencoded, got", result.Header.Get("Content-Type"))
	}

	_, err = newRequest("://example.com/api", nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestKeyFromUrl(t *testing.T) {
	result := keyFromURL("https://pastebin.com/pasteKey")
	if result != "pasteKey" {
		t.Error("Expected pasteKey, got", result)
	}
	result = keyFromURL("")
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}
