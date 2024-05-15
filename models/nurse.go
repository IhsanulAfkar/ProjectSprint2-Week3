package models

type Nurse struct {
	Id                   string  `db:"id" json:"userId"`
	Nip                  int64   `db:"nip" json:"nip"`
	Name                 string  `db:"name" json:"name"`
	Password             *string `db:"password" json:"password,omitempty"`
	IsGranted            bool    `db:"isGranted" json:"isGranted"`
	IdentityCardScanning string  `db:"identityCardScanning" json:"identityCardScanning"`
	CreatedAt            string  `db:"createdAt" json:"createdAt"`
	UpdatedAt            string  `db:"updatedAt" json:"updatedAt"`
	DeletedAt            *string `db:"deletedAt" json:"deletedAt,omitempty"`
}