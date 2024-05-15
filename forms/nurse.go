package forms

type NurseRegister struct {
	Nip                 int64  `json:"nip"`
	Name                string `json:"name"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}
type NurseUpdate struct {
	Nip  int64  `json:"nip"`
	Name string `json:"name"`
}