package pastebin

import (
	"fmt"
	"time"
)

type Visibility int

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

func (p *Paste) String() string {
	return fmt.Sprintf("Key: %s, Title: %s, URL: %s, CreatedAt: %s, ExpireDate: %s, Visibility: %d, FormatLong: %s",
		p.Key, p.Title, p.URL, p.CreatedAt.Format(time.RFC3339), p.ExpireDate.Format(time.RFC3339), p.Visibility, p.FormatLong)
}

type pastesXML struct {
	Pastes []pasteXML `xml:"paste"`
}

type pasteXML struct {
	Key         string `xml:"paste_key"`
	Date        int64  `xml:"paste_date"`
	Title       string `xml:"paste_title"`
	Size        int    `xml:"paste_size"`
	ExpireDate  int64  `xml:"paste_expire_date"`
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
		CreatedAt:   time.Unix(p.Date, 0),
		ExpireDate:  time.Unix(p.ExpireDate, 0),
		Visibility:  Visibility(p.Private),
		FormatLong:  p.FormatLong,
		FormatShort: p.FormatShort,
	}
}
