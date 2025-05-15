package pastebin

import (
	"testing"
	"time"
)

const defaultPasteStringResult = "Key: pastekey, Title: Test, URL: https://pastebin.com/pastekey, CreatedAt: 2023-10-01T12:00:00Z, ExpireDate: 2023-10-31T12:00:00Z, Visibility: Private, FormatLong: Go"

func TestVisibilityTypeString(t *testing.T) {
	publicExp := "Public"
	unlistedExp := "Unlisted"
	privateExp := "Private"
	unknownExp := "Unknown"
	if Public.String() != publicExp {
		t.Errorf("Expected: %s\nGot: %s", publicExp, Public.String())
	}
	if Unlisted.String() != unlistedExp {
		t.Errorf("Expected: %s\nGot: %s", unlistedExp, Unlisted.String())
	}
	if Private.String() != privateExp {
		t.Errorf("Expected: %s\nGot: %s", privateExp, Private.String())
	}
	if Visibility(3).String() != unknownExp {
		t.Errorf("Expected: %s\nGot: %s", unknownExp, Visibility(3).String())
	}
}

func TestPasteString(t *testing.T) {
	paste := &Paste{
		Key:        "pastekey",
		Title:      "Test",
		URL:        "https://pastebin.com/pastekey",
		CreatedAt:  time.Unix(1696161600, 0).UTC(),
		ExpireDate: time.Unix(1698753600, 0).UTC(),
		Visibility: 2,
		FormatLong: "Go",
	}
	if paste.String() != defaultPasteStringResult {
		t.Errorf("Expected: %s\nGot: %s", defaultPasteStringResult, paste.String())
	}
}

func TestPasteXMLToPaste(t *testing.T) {
	expected := defaultPasteStringResult
	paste := pasteXML{
		Key:        "pastekey",
		Title:      "Test",
		Date:       1696161600,
		ExpireDate: 1698753600,
		Private:    2,
		FormatLong: "Go",
		URL:        "https://pastebin.com/pastekey",
	}.toPaste()
	if expected != paste.String() {
		t.Errorf("Expected: %s\nGot: %s", expected, paste.String())
	}
}
