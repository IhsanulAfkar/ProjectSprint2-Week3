package forms

type AdminRegister struct {
	Nip      int64  `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type AdminLogin struct {
	Nip      int64  `json:"nip"`
	Password string `json:"password"`
}