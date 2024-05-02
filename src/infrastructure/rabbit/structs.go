package rabbit

type EventUser struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}
