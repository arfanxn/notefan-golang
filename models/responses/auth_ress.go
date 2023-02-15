package responses

type AuthLogin struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type AuthRegister struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
