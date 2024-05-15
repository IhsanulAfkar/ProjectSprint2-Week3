package forms

type PatientCreate struct {
	IdentityNumber      int64  `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}