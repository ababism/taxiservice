package models

type Error struct {
	Message string `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}
