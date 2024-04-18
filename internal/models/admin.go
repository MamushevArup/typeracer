package models

type AdminSignIn struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type AdminSignInResponse struct {
	Access string `json:"access"`
}
