package models

type AuthResponse struct {
	Access string `json:"access"`
}
type SignIn struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type SignUp struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type RefreshS struct {
	Fingerprint string `json:"fingerprint"`
}
