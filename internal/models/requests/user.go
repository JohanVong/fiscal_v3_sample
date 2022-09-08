package requests

type CreateUserRequestMerchant struct {
	IdCompany   int    `json:"IdCompany" valid:"required~Поле IdCompany не должно быть пустым"`
	PhoneLogin  string `json:"PhoneLogin" valid:"required~Поле PhoneLogin не должно быть пустым,isPhoneNumber~Не является номером телефона"`
	PasswordRaw string `json:"Password" valid:"required~Поле Password не должно быть пустым,stringlength(8|32)~Пароль должен быть длиной от 8 до 32 символов"`
	Name        string `json:"Name"`
	Groups      []int  `json:"Groups"`
}

type CreateUserRequestSuperAdmin struct {
	IdCompany   int    `json:"IdCompany"`
	PhoneLogin  string `json:"PhoneLogin" valid:"required~Поле PhoneLogin не должно быть пустым,isPhoneNumber~Не является номером телефона"`
	PasswordRaw string `json:"Password" valid:"required~Поле Password не должно быть пустым,stringlength(8|32)~Пароль должен быть длиной от 8 до 32 символов"`
	Name        string `json:"Name"`
	Groups      []int  `json:"Groups"`
}

type EditUserRequest struct {
	IdTypeUser  int    `json:"IdTypeUser"`
	PhoneLogin  string `json:"PhoneLogin" valid:"isPhoneNumber~Не является номером телефона"`
	PasswordRaw string `json:"Password" valid:"stringlength(8|32)~Пароль должен быть длиной от 8 до 32 символов"`
	Name        string `json:"Name"`
}
