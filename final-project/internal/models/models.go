package models

type Album struct {
	ID     string `json:"id,omitempty"`
	Artist string `json:"artist,omitempty"`
	Title  string `json:"title,omitempty"`
	Year   int    `json:"year,omitempty"`
}

type User struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Object struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Review struct {
	ID       int    `json:"id,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	ObjectID int    `json:"object_id,omitempty"`
	Text     string `json:"text,omitempty,"`
	Rating   int    `json:"rating,omitempty"`
}

type ReviewRepresentation struct {
	ID       int    `json:"id,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	ObjectID int    `json:"object_id,omitempty"`
	Text     string `json:"text,omitempty,"`
	Rating   string `json:"rating,omitempty"`
}
