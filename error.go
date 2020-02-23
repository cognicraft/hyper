package hyper

// Error .
type Error struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	Message     string `json:"message"`
	Code        string `json:"code,omitempty"`
}

// Errors .
type Errors []Error
