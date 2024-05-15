package models

type Record struct {
	Id             string `db:"id" json:"userId"`
	IdentityNumber int64  `db:"identityNumber" json:"identityNumber"`
	Nip            int64  `db:"nip" json:"nip"`
	Symptoms       string `db:"symptoms" json:"symptoms"`
	Medications    string `db:"medications" json:"medications"`
	CreatedAt      string `db:"createdAt" json:"createdAt"`
}