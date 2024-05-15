package models

var Role = [2]string{
	"it",
	"nurse",
}

type Admin struct {
	Id        string `db:"id" json:"userId"`
	Nip       int64  `db:"nip" json:"nip"`
	Name      string `db:"name" json:"name"`
	Password  string `db:"password" json:"password"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
}

type NIP struct {
	ToString     string
	ToInt        int64
	First3Digits string
	Gender       string
	Year         string
	Month        string
	EndDigits    string
}

type GetUser struct {
	Id        string `db:"id" json:"userId"`
	Nip       int64  `db:"nip" json:"nip"`
	Name      string `db:"name" json:"name"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
}