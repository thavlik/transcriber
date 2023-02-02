package definitions

type RemoteIAM interface {
	DeleteUser(DeleteUser) Void
	Login(LoginRequest) LoginResponse
}

type DeleteUser struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	DeleteProjects bool   `json:"deleteProjects"`
}

type Void struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
