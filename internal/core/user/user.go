package user

type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Admin     *bool   `json:"admin"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}
