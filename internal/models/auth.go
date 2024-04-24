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

type SignUpService struct {
	Access  string
	Avatar  string
	Refresh string
}

type SignUpHandler struct {
	Access string `json:"access"`
	Avatar string `json:"avatar"`
}

type SignInService struct {
	Access   string
	Refresh  string
	Username string
	Avatar   string
}

type SignInHandler struct {
	Access   string `json:"access"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
