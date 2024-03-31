package dtos

type ActionUserDto struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type ChatActionDto struct {
	Action      string          `json:"action"`
	ActionUsers []ActionUserDto `json:"action_users"`
}

type ChatDto struct {
	Id         int              `json:"id"`
	Avatar     FileDto          `json:"avatar"`
	Title      string           `json:"title"`
	Type       string           `json:"type"`
	Members    []int            `json:"members"`
	IsArchived bool             `json:"is_archived"`
	OwnerId    int              `json:"owner_id"`
	Admins     []int            `json:"admins"`
	Actions    *[]ChatActionDto `json:"actions"`
}
