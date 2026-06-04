package slaughtergroups

import "time"

type GroupDB struct {
	Id                string    `db:"id"`
	EntrytDate        time.Time `db:"entry_date"`
	ButcherId         string    `db:"butcher_id"`
	ButcherName       string    `db:"butcher_name"`
	ButcherDiscount   *float64  `db:"butcher_discount"`
	Discount          float64   `db:"discount"`
	AnimalsNumber     *int      `db:"animals_number"`
	AverageWeight     *float64  `db:"average_weight"`
	AverageDeadWeight *float64  `db:"average_dead_weight"`
	AverageRate       *float64  `db:"average_rate"`
	CreatedAt         time.Time `db:"created_at"`
}

type GroupSave struct {
	Id           *string   `db:"id"`
	EntrytDate   time.Time `db:"entry_date"`
	ButcherId    string    `db:"butcher_id"`
	DiscountRate float64   `db:"discount_rate"`
}
