package order

type OrderReq struct {
	CouponCode *string      `json:"couponCode,omitempty"`
	Items      *[]OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type Order struct {
	ID         string       `json:"id"`
	Items      []OrderItem  `json:"items"`
	Products   []ProductRef `json:"products"`
	CouponCode *string      `json:"couponCode,omitempty"`
}

type ProductRef struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}
