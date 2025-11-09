package model

type CreateCourierRequest struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type UpdateCourierRequest struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}