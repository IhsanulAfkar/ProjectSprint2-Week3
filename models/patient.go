package models

var Gender = [2]string{
	"male",
	"female",
}

type Patient struct {
	Id                   string `db:"id" json:"userId"`
	IdentityNumber       int64  `db:"identityNumber" json:"identityNumber"`
	Name                 string `db:"name" json:"name"`
	PhoneNumber          string `db:"phoneNumber" json:"phoneNumber"`
	BirthDate            string `db:"birthDate" json:"birthDate"`
	Gender               string `db:"gender" json:"gender"`
	IdentityCardScanning string `db:"identityCardScanning" json:"identityCardScanning"`
	CreatedAt            string `db:"createdAt" json:"createdAt"`
}