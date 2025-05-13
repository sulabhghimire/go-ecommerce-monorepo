package dto

type CreateCartRequest struct {
	ProductId uint `json:"product_id"`
	Quantity  uint `json:"qty"`
}
