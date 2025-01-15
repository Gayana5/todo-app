package todo

type TodoList struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Descriprtion string `json:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}
type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}
