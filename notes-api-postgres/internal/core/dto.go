package core

type NoteCreate struct {
	Title   string `json:"title" example:"Новая заметка"`
	Content string `json:"content" example:"Текст заметки"`
}

type NoteUpdate struct {
	Title   *string `json:"title,omitempty" example:"Обновлено"`
	Content *string `json:"content,omitempty" example:"Новый текст"`
}