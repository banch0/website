package dto

// BurgerDTO Data Transfer Object
type BurgerDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
