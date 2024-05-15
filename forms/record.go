package forms

type CreateRecord struct {
	IdentityNumber int64  `json:"identityNumber"`
	Symptoms       string `json:"symptoms"`
	Medications    string `json:"medications"`
}