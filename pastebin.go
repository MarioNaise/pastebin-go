package pastebin

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Expiration defines the duration before a paste expires.
type Expiration string

// TODO: implement NewCreatePasteRequest func
// TODO: improve error messages

// CreatePasteRequest holds the parameters to create a new paste.
//
// See https://pastebin.com/doc_api#2
type CreatePasteRequest struct {
	// required.
	// this is the text that will be written inside your paste.
	Content string

	// optional.
	// this will be the name / title of your paste.
	Name string

	// optional.
	// this will be the syntax highlighting value.
	//
	// See https://pastebin.com/doc_api#5
	Format string

	// optional.
	// this sets the key of the folder of your paste.
	//
	// See https://pastebin.com/doc_api#5
	Folder string

	// optional.
	// this sets the expiration date of your paste.
	// default value: "N" (Never)
	//
	// See https://pastebin.com/doc_api#6
	Expiration Expiration

	// optional.
	// this makes a paste public, unlisted or private.
	// Public = 0, Unlisted = 1, Private = 2
	//
	// See https://pastebin.com/doc_api#7
	Visibility Visibility

	// optional.
	// if true, this will create the paste as the currently logged in user.
	// otherwise it will create the paste as a guest.
	CreatePasteAsUser bool
}

// Client is the Pastebin API client.
type Client struct {
	apiUserName     string
	apiUserPassword string
	apiDevKey       string
	apiUserKey      string
}

// NewClient creates a new Pastebin API client.
// If a username is provided, it logs the user in to obtain a user API key.
//
// See https://pastebin.com/doc_api#9
func NewClient(userName, password, devKey string) (*Client, error) {
	client := &Client{
		apiUserName:     userName,
		apiUserPassword: password,
		apiDevKey:       devKey,
	}
	if len(userName) > 0 {
		return client, client.refreshUserKey()
	}
	return client, nil
}

// CreatePaste creates a new paste using the given request parameters.
func (c *Client) CreatePaste(req *CreatePasteRequest) (string, error) {
	vals := url.Values{
		apiDevKey:    {c.apiDevKey},
		apiPasteCode: {req.Content},
		apiOption:    {"paste"},
	}
	params := [4][2]string{
		{apiPasteName, req.Name},
		{apiPasteFormat, req.Format},
		{apiPasteExpireDate, string(req.Expiration)},
		{apiFolderKey, req.Folder},
	}
	for _, pair := range params {
		if len(pair[1]) > 0 {
			vals.Add(pair[0], pair[1])
		}
	}
	if req.CreatePasteAsUser {
		vals.Add(apiUserKey, c.apiUserKey)
	}
	if req.Visibility > 0 {
		vals.Add(apiPastePrivate, fmt.Sprintf("%d", req.Visibility))
	}
	resp, err := c.post(PostUrl, vals, true)
	if err != nil {
		return "", err
	}
	return keyFromURL(resp), nil
}

// GetUserPastes retrieves the list of pastes created by the authenticated user.
func (c *Client) GetUserPastes() ([]*Paste, error) {
	resp, err := c.post(PostUrl, url.Values{
		apiDevKey:       {c.apiDevKey},
		apiUserKey:      {c.apiUserKey},
		apiResultsLimit: {"250"},
		apiOption:       {"list"},
	}, true)
	if err != nil {
		return nil, err
	}
	var pastes pastesXML
	if err := xml.Unmarshal(fmt.Appendf(nil, "<root>%s</root>", resp), &pastes); err != nil {
		return nil, err
	}
	p := make([]*Paste, len(pastes.Pastes))
	for i, pasteXML := range pastes.Pastes {
		p[i] = pasteXML.toPaste()
	}
	return p, nil
}

// DeletePaste deletes a paste by its unique key.
func (c *Client) DeletePaste(key string) error {
	res, err := c.post(PostUrl, url.Values{
		apiDevKey:   {c.apiDevKey},
		apiUserKey:  {c.apiUserKey},
		apiOption:   {"delete"},
		apiPasteKey: {key},
	}, true)
	if err != nil {
		return err
	}
	if res != "Paste Removed" {
		return errors.New("unexpected Pastebin response: " + res)
	}
	return nil
}

// GetRawUserPasteContent retrieves the raw content of a user-owned paste.
func (c *Client) GetRawUserPasteContent(key string) (string, error) {
	resp, err := c.post(RawUrl, url.Values{
		apiDevKey:   {c.apiDevKey},
		apiUserKey:  {c.apiUserKey},
		apiPasteKey: {key},
		apiOption:   {"show_paste"},
	}, true)
	if err != nil {
		return "", err
	}
	return resp, nil
}

// GetRawPublicPasteContent fetches the raw content of a public or unlisted paste.
func (c *Client) GetRawPublicPasteContent(key string) (string, error) {
	return c.get(RawPublicUrl+"/"+key, url.Values{}, true)
}

// GetUserDetails retrieves account details of the authenticated user.
func (c *Client) GetUserDetails() (*User, error) {
	resp, err := c.post(PostUrl, url.Values{
		apiDevKey:  {c.apiDevKey},
		apiUserKey: {c.apiUserKey},
		apiOption:  {"userdetails"},
	}, true)
	if err != nil {
		return nil, err
	}
	var user User
	if err := xml.Unmarshal([]byte(resp), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) refreshUserKey() error {
	resp, err := c.post(LoginUrl, url.Values{
		apiUserName:     {c.apiUserName},
		apiUserPassword: {c.apiUserPassword},
		apiDevKey:       {c.apiDevKey},
	}, false)
	if err != nil {
		return err
	}
	c.apiUserKey = resp
	return nil
}

func (c *Client) get(u string, _ url.Values, reauth bool) (string, error) {
	resp, err := getHttpClient().Get(u)
	if err != nil {
		return "", err
	}
	return c.handleResponse(resp, c.get, u, url.Values{}, reauth)
}

func (c *Client) post(url string, vals url.Values, reauth bool) (string, error) {
	req, err := newRequest(url, vals)
	if err != nil {
		return "", err
	}
	resp, err := getHttpClient().Do(req)
	if err != nil {
		return "", err
	}
	return c.handleResponse(resp, c.post, url, vals, reauth)
}

func (c *Client) handleResponse(resp *http.Response,
	cb func(string, url.Values, bool) (string, error),
	url string, vals url.Values, reauth bool,
) (string, error) {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	res := string(body)
	if res == invalidApiUserKeyResponse &&
		resp.StatusCode != 200 && reauth {
		err := c.refreshUserKey()
		if err != nil {
			return "", err
		}
		return cb(url, vals, false)
	}
	if resp.StatusCode != 200 {
		return "", errors.New(res)
	}
	return res, nil
}

func newRequest(url string, vals url.Values) (*http.Request, error) {
	reqBody := strings.NewReader(vals.Encode())
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func keyFromURL(url string) string {
	return strings.TrimPrefix(url, BaseUrl+"/")
}
