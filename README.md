# Pastebin Go Client

- [Pastebin Go Client](#pastebin-go-client)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
  - [Usage](#usage)
    - [Importing](#importing)
    - [Initialize Client](#initialize-client)
    - [Create a Paste](#create-a-paste)
    - [Get User Pastes](#get-user-pastes)
    - [Delete a Paste](#delete-a-paste)
    - [Get Raw Paste Content](#get-raw-paste-content)
      - [For User Authenticated Pastes](#for-user-authenticated-pastes)
      - [For Public or Unlisted Pastes](#for-public-or-unlisted-pastes)
    - [Get User Details](#get-user-details)
  - [Visibility Options](#visibility-options)
  - [Expiration Options](#expiration-options)
  - [License](#license)

This is a Go client library for interacting with the [Pastebin API](https://pastebin.com/doc_api). It allows users to authenticate, create pastes, delete them, access raw content and retrieve user pastes.

---

## Features

- Authenticate using Pastebin credentials
- Create pastes (public, unlisted, or private)
- Fetch all pastes for a user
- Delete a paste
- Retrieve raw content from a user's or public paste
- Fetch user details

---

## Getting Started

### Installation

```bash
go get github.com/MarioNaise/pastebin-go
```

---

## Usage

### Importing

```go
import "github.com/MarioNaise/pastebin-go"
```

---

### Initialize Client

```go
client, err := pastebin.NewClient("your_username", "your_password", "your_dev_key")
if err != nil {
 log.Fatal(err)
}
```

---

### Create a Paste

```go
pasteKey, err := client.CreatePaste(&pastebin.CreatePasteRequest{
  Content:           "Hello, Pastebin!", // required
  Name:              "HelloWorld",       // optional
  Format:            "text",             // optional
  Folder:            "hello",            // optional
  Expiration:        pastebin.OneHour,   // optional
  Visibility:        pastebin.Unlisted,  // optional
  CreatePasteAsUser: true,               // optional
})
if err != nil {
 log.Fatal(err)
}
fmt.Println("Paste Key:", pasteKey)
```

---

### Get User Pastes

```go
pastes, err := client.GetUserPastes()
if err != nil {
 log.Fatal(err)
}
for _, p := range pastes {
 fmt.Println(p)
}
```

---

### Delete a Paste

```go
err := client.DeletePaste("paste_key")
if err != nil {
 log.Fatal(err)
}
```

---

### Get Raw Paste Content

#### For User Authenticated Pastes

```go
content, err := client.GetRawUserPasteContent("paste_key")
if err != nil {
 log.Fatal(err)
}
fmt.Println(content)
```

#### For Public or Unlisted Pastes

```go
content, err := client.GetRawPublicPasteContent("public_or_unlisted_paste_key")
if err != nil {
 log.Fatal(err)
}
fmt.Println(content)
```

---

### Get User Details

```go
user, err := client.GetUserDetails()
if err != nil {
 log.Fatal(err)
}
fmt.Printf("User: %s, Type: %s\n", user.UserName, user.AccountType)
```

---

## Visibility Options

```go
pastebin.Public   // 0
pastebin.Unlisted // 1
pastebin.Private  // 2
```

---

## Expiration Options

```go
pastebin.Never
pastebin.TenMinutes
pastebin.OneHour
pastebin.OneDay
pastebin.OneWeek
pastebin.TwoWeeks
pastebin.OneMonth
pastebin.SixMonths
pastebin.OneYear
```

---

## License

MIT License
