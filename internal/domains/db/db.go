package db

import "github.com/google/uuid"

type User struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Country *string   `json:"country,omitempty"`
	Age     int32     `json:"age"`
}

type ID struct {
	Id uuid.UUID `json:"id"`
}

type CreateRequest struct {
	Name    string  `json:"name"`
	Surname string  `json:"surname"`
	Country *string `json:"country,omitempty"`
	Age     int32   `json:"age"`
}

type UpdateRequest struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
}

type DeleteRequest struct {
	Id uuid.UUID `json:"id"`
}
