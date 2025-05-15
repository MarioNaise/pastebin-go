package pastebin

import "testing"

func TestAccountTypeString(t *testing.T) {
	normalExpected := "NormalUser"
	proExpected := "ProUser"
	unknownExpected := "UnknownUserType"
	if NormalUser.String() != normalExpected {
		t.Errorf("Expected: %s\nGot: %s", normalExpected, NormalUser.String())
	}
	if ProUser.String() != proExpected {
		t.Errorf("Expected: %s\nGot: %s", proExpected, ProUser.String())
	}
	if AccountType(2).String() != unknownExpected {
		t.Errorf("Expected: %s\nGot: %s", unknownExpected, AccountType(2).String())
	}
}

func TestUserString(t *testing.T) {
	user := User{
		UserName: "test", Expiration: "N",
		Visibility: Public, Avatar: "http://example.com/test.png",
		Website: "http://example.com", Email: "test@example.com",
		Location: "Test-Location", AccountType: NormalUser,
	}
	expected := "UserName: test, Expiration: N, Visibility: Public, Avatar: http://example.com/test.png, Website: http://example.com, Email: test@example.com, Location: Test-Location, AccountType: NormalUser"
	if expected != user.String() {
		t.Errorf("Expected: %s\nGot: %s", expected, user.String())
	}
}
