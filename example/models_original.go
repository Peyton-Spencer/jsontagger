package example

// Example struct with mixed case JSON tags
type User struct {
	ID        int    `json:"user_id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Address   Address `json:"address"`
}

type Address struct {
	Street  string `json:"street_name"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode,omitempty"`
}

type Product struct {
	ProductID   int     `json:"productId"`
	Name        string  `json:"name"`
	Description string  `json:"product_description"`
	Price       float64 `json:"price,omitempty"`
	Categories  []string `json:"categories"`
}