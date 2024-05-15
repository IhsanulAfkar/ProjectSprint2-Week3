package models

type Image struct {
	Id        string `db:"id" json:"id"`
	Path      string `db:"path" json:"path"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
}