package structs

type UserCommint struct {
	Username string `josn:"username"`
	Commint  string `json:"Commint"`
	Stars    int    `json:"stars"`
}

type Color struct {
	ColorName string `json:"color"`
	Qty       int    `json:"qty"`
}

type Size struct {
	InSotck bool    `json:"inSotck"`
	Colors  []Color `json:"colors"`
}

type Model struct {
	Id          int             `json:"id"`
	Sizes       map[string]Size `json:"sizes"`
	Price       int             `json:"price"`
	Description string          `json:"description"`
	Discount    int             `json:"discount"`
	LinkesImage []string        `json:"linkesImage"`
	Commint     []UserCommint   `json:"commint"`
}
