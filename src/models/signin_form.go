package models

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Host     string `json:"host"`
	Port     int	`json:"port"`
}

type SignInFormValid struct {
	HasError bool   `json:"has_error"`
	Message  interface{} `json:"message"`
	Addition interface{} `json:"addition"`
}
