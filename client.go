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

type Client struct {
	apiUserName     string
	apiUserPassword string
	apiDevKey       string
	apiUserKey      string
}

func NewClient(userName, password, devKey string) (*Client, error) {
	client := &Client{
		apiUserName:     userName,
		apiUserPassword: password,
		apiDevKey:       devKey,
	}
	if len(userName) > 0 {
		return client, client.login()
	}
	return client, nil
}

func (c *Client) CreatePaste(req *CreatePasteRequest) (string, error) {
	vals := url.Values{
		apiDevKey:    {c.apiDevKey},
		apiPasteCode: {req.Content},
		apiOption:    {"paste"},
	}
	if req.CreatePasteAsUser {
		vals.Add(apiUserKey, c.apiUserKey)
	}
	if len(req.Name) > 0 {
		vals.Add(apiPasteName, req.Name)
	}
	if len(req.Format) > 0 {
		vals.Add(apiPasteFormat, req.Format)
	}
	if len(req.Expiration) > 0 {
		vals.Add(apiPasteExpireDate, string(req.Expiration))
	}
	if len(req.Folder) > 0 {
		vals.Add(apiFolderKey, req.Folder)
	}
	if req.Visibility > 0 {
		vals.Add(apiPastePrivate, fmt.Sprintf("%d", req.Visibility))
	}
	resp, err := c.do(PostUrl, vals, true)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *Client) GetUserPastes() ([]*Paste, error) {
	resp, err := c.do(PostUrl, url.Values{
		apiDevKey:       {c.apiDevKey},
		apiUserKey:      {c.apiUserKey},
		apiResultsLimit: {"250"},
		apiOption:       {"list"},
	}, true)
	if err != nil {
		return []*Paste{}, err
	}
	var pastes pastesXML
	if err := xml.Unmarshal(fmt.Appendf(nil, "<root>%s</root>", resp), &pastes); err != nil {
		return []*Paste{}, err
	}
	p := make([]*Paste, len(pastes.Pastes))
	for i, paste := range pastes.Pastes {
		p[i] = paste.toPaste(c.apiUserName)
	}
	return p, nil
}

func (c *Client) DeletePaste(key string) error {
	_, err := c.do(PostUrl, url.Values{
		apiDevKey:   {c.apiDevKey},
		apiUserKey:  {c.apiUserKey},
		apiOption:   {"delete"},
		apiPasteKey: {key},
	}, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetRawUserPasteContent(key string) (string, error) {
	resp, err := c.do(RawUrl, url.Values{
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

func (c *Client) GetRawPublicPasteContent(key string) (string, error) {
	resp, err := getHttpClient().Get(RawPublicUrl + "/" + key)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) GetUserDetails() (*User, error) {
	resp, err := c.do(PostUrl, url.Values{
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

func (c *Client) login() error {
	return c.refreshUserKey()
}

func (c *Client) refreshUserKey() error {
	resp, err := c.do(LoginUrl, url.Values{
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

func (c *Client) do(url string, vals url.Values, reauthenticateOnError bool) (string, error) {
	req, err := newRequest(url, vals)
	if err != nil {
		return "", err
	}
	resp, err := getHttpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := string(body)
	if response == "Bad API request, invalid api_user_key" &&
		reauthenticateOnError {
		err := c.refreshUserKey()
		if err != nil {
			return "", err
		}
		return c.do(url, vals, false)
	}
	if strings.HasPrefix(response, "Bad API request") {
		return "", errors.New(response)
	}
	return response, nil
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
