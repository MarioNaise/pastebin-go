package pastebin

type Expiration string

type CreatePasteRequest struct {
	CreatePasteAsUser bool
	Content           string
	Name              string
	Format            string
	Expiration        Expiration
	Folder            string
	Visibility        Visibility
}
