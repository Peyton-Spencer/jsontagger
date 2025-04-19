package example

// Struct without any JSON tags
type Person struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
}

// Struct with mixed tags (some fields have tags, some don't)
type Company struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
	Founded int    `json:"founded"`
}

// Struct with non-JSON tags
type Project struct {
	ID          int    `db:"project_id" json:"id"`
	Name        string `validate:"required" json:"name"`
	Description string `json:"description"`
}
