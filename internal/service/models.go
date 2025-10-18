package service

import "time"

// Entidades do Domínio
type Author struct {
	ID   int    `json:"-" db:"id"`    // PK interna, não exposta no JSON
	UUID string `json:"id" db:"uuid"` // ID público
	Name string `json:"name" db:"name"`
}

type Category struct {
	ID   int    `json:"-" db:"id"`
	UUID string `json:"id" db:"uuid"`
	Name string `json:"name" db:"name"`
}

type Book struct {
	ID          int       `json:"-" db:"id"`    // PK interna
	UUID        string    `json:"id" db:"uuid"` // ID público (UUID)
	Title       string    `json:"title" db:"title"`
	Author      Author    `json:"author" db:"author"`
	Category    Category  `json:"category" db:"category"`
	PublishedAt int       `json:"published_at" db:"published_at"`
	ISBN        string    `json:"isbn" db:"isbn"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Struct de Entrada para o método de Serviço (POST /books)
type CreateBookInput struct {
	Title            string            `json:"title"`
	AuthorName       string            `json:"author_name"`   // Nome do autor para criar/vincular
	CategoryName     string            `json:"category_name"` // Nome da categoria para criar/vincular
	PublishedAt      int               `json:"published_at"`
	ISBN             string            `json:"isbn"`
	ValidationErrors map[string]string `json:"-"` // Para validação
}
