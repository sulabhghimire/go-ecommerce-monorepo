package dto

type SellerOrderDetails struct {
	OrderRefNumber  uint    `json:"order_ref_number"`
	OrderStatus     string  `json:"order_status"`
	CreatedAt       string  `json:"created_at"`
	OrderItemId     uint    `json:"order_item_id"`
	ProductId       uint    `json:"product_id"`
	Name            string  `json:"name"`
	ImageUrl        string  `json:"image_url"`
	Price           float64 `json:"price"`
	Qty             uint    `json:"qty"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email"`
	CustomerAddress string  `json:"customer_address"`
}
