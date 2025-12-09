package slaughter

type ButcherEntry struct {
	Id            string   `json:"id" db:"id"`
	Name          string   `json:"name" db:"name"`
	Cnpj          *string  `json:"cnpj" db:"cnpj"`
	Address       *string  `json:"address" db:"address"`
	Discount      *float64 `json:"discount" db:"discount"`
	AnimalsNumber int      `json:"animalsNumber" db:"animals_number"`
	AverageWeight float64  `json:"averageWeight" db:"avg_weight"`
	AverageRate   float64  `json:"averageRate" db:"avg_rate"`
}

type ButcherSave struct {
	Id       string   `json:"id" db:"id"`
	Name     string   `json:"name" db:"name"`
	Cnpj     *string  `json:"cnpj" db:"cnpj"`
	Address  *string  `json:"address" db:"address"`
	Discount *float64 `json:"discount" db:"discount"`
	UserId   string   `json:"-" db:"user_id"`
}
