package pastebin

import (
	"fmt"
	"time"
)

// Visibility defines the visibility level of a paste.
// Public = 0, Unlisted = 1, Private = 2.
//
// See https://pastebin.com/doc_api#7
type Visibility int

// String returns the string representation of a Visibility.
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

// Paste represents a Pastebin paste entry.
type Paste struct {
	Key         string
	Title       string
	URL         string
	Hits        int
	Size        int
	CreatedAt   time.Time
	ExpireDate  time.Time
	Visibility  Visibility
	FormatLong  string
	FormatShort string
}

// String returns a formatted string of the paste data.
func (p *Paste) String() string {
	return fmt.Sprintf("Key: %s, Title: %s, URL: %s, CreatedAt: %s, ExpireDate: %s, Visibility: %s, FormatLong: %s",
		p.Key, p.Title, p.URL, p.CreatedAt.Format(time.RFC3339), p.ExpireDate.Format(time.RFC3339), p.Visibility, p.FormatLong)
}

type pastesXML struct {
	Pastes []pasteXML `xml:"paste"`
}

type pasteXML struct {
	Key         string `xml:"paste_key"`
	Title       string `xml:"paste_title"`
	Date        int64  `xml:"paste_date"`
	ExpireDate  int64  `xml:"paste_expire_date"`
	Size        int    `xml:"paste_size"`
	Private     int    `xml:"paste_private"`
	FormatLong  string `xml:"paste_format_long"`
	FormatShort string `xml:"paste_format_short"`
	URL         string `xml:"paste_url"`
	Hits        int    `xml:"paste_hits"`
}

func (p pasteXML) toPaste() *Paste {
	return &Paste{
		Key:         p.Key,
		Title:       p.Title,
		URL:         p.URL,
		Hits:        p.Hits,
		Size:        p.Size,
		CreatedAt:   time.Unix(p.Date, 0).UTC(),
		ExpireDate:  time.Unix(p.ExpireDate, 0).UTC(),
		Visibility:  Visibility(p.Private),
		FormatLong:  p.FormatLong,
		FormatShort: p.FormatShort,
	}
}
