package requests

type AuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=2,max=50"`
}

type AuthRegister struct {
	Name            string `json:"name" validate:"required,min=2,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=50"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}
