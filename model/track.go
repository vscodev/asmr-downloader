package model

type Track struct {
	Type             string   `json:"type"`
	Title            string   `json:"title"`
	MediaDownloadURL string   `json:"mediaDownloadUrl"`
	Children         []*Track `json:"children"`
}

func (t Track) IsFolder() bool {
	return t.Type == "folder"
}

func (t Track) String() string {
	return t.Title
}
