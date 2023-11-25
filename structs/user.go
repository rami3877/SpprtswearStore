package structs

type Email struct {
	IsActive  bool   `json:"isActive"`
	EmailName string `json:"emailName"`
}
type Addr struct {
	Address   bool   `json:"address"`
	City      string `json:"city"`
	Street    string `json:"street"`
	Apartment int    `json:"apartment"`
	Building  int    `json:"building"`
}
type Visa struct {
	Number         string `json:"number"`
	Cvv            string `json:"cvv"`
	ExpirationData string `json:"expirationData"`
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserEmail Email  `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	UserAddr  Addr   `json:"addr"`
	UserVisa  []Visa `json:"visa"`
}
