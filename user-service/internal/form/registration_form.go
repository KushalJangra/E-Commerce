package form

type RegistrationForm struct {
	FirstName string `json:"FirstName" binding:"required"`
	LastName  string `json:"LastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	Email     string `json:"email" binding:"required,email"`
}
