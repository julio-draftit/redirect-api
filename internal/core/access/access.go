package access

type Access struct {
	ID         int     `json:"id"`
	RedirectID int     `json:"redirect_id"`
	CreatedAt  *string `json:"created_at"`
}
