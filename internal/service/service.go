package service

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrFailedValidation = errors.New("failed validation")

var ErrRecordNotFound = errors.New("record not found")

// O Fat Service Struct que armazena as dependências (DB)
type Service struct {
	DB *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		DB: db,
	}
}

// Método principal para criar um livro - Contém toda a lógica de negócio
func (s *Service) CreateBook(input *CreateBookInput) (*Book, error) {
	// 1. Validação de Entrada
	input.ValidationErrors = make(map[string]string)
	if input.Title == "" {
		input.ValidationErrors["title"] = "must be provided"
	}
	if input.AuthorName == "" {
		input.ValidationErrors["author_name"] = "must be provided"
	}
	// ... (Outras validações)

	if len(input.ValidationErrors) > 0 {
		return nil, ErrFailedValidation
	}

	// 2. Lógica de Negócios e DB

	// Iniciamos uma transação para garantir atomicidade (Autor, Categoria e Livro)
	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Exemplo: Lógica para garantir que a Categoria exista (ou criar)
	categoryID, categoryUUID, err := s.getOrCreateCategory(tx, input.CategoryName)
	if err != nil {
		return nil, err
	}

	// Exemplo: Lógica para garantir que o Autor exista (ou criar)
	authorID, authorUUID, err := s.getOrCreateAuthor(tx, input.AuthorName)
	if err != nil {
		return nil, err
	}

	bookUUID := uuid.New().String()

	query := `
        INSERT INTO books (uuid, title, author_id, category_id, published_at, isbn, created_at) 
        VALUES (?, ?, ?, ?, ?, ?, datetime('now')) 
        RETURNING id, created_at
    `
	var newBook Book
	err = tx.QueryRowx(query,
		bookUUID, input.Title, authorID, categoryID, input.PublishedAt, input.ISBN,
	).Scan(&newBook.ID, &newBook.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 4. Popular o struct (usando os UUIDs para a resposta)
	newBook.UUID = bookUUID // O ID público
	newBook.Title = input.Title
	newBook.PublishedAt = input.PublishedAt
	newBook.ISBN = input.ISBN

	// Usamos os UUIDs (públicos) nos structs aninhados
	newBook.Author = Author{UUID: authorUUID, Name: input.AuthorName}
	newBook.Category = Category{UUID: categoryUUID, Name: input.CategoryName}

	return &newBook, nil
}

// Funções auxiliares (internas ao serviço, mantêm a lógica encapsulada)

func (s *Service) getOrCreateCategory(tx *sqlx.Tx, name string) (int, string, error) { // Retorna ID interno E UUID
	var id int
	var currentUUID string
	query := "SELECT id, uuid FROM categories WHERE name = ?"
	err := tx.QueryRowx(query, name).Scan(&id, &currentUUID)
	if err == nil {
		// Categoria encontrada: retornamos os dados dela.
		return id, currentUUID, nil
	}

	// Verificamos se o erro foi apenas que a linha não existe.
	if errors.Is(err, sql.ErrNoRows) {
		// Categoria não existe, vamos criar:

		newUUID := uuid.New().String() // Gerar NOVO UUID

		// 2. Inserir um novo registro, incluindo o UUID.
		result, err := tx.Exec("INSERT INTO categories (name, uuid) VALUES (?, ?)", name, newUUID)
		if err != nil {
			return 0, "", err
		}

		lastID, _ := result.LastInsertId()

		// 3. Retornar os IDs do novo registro.
		return int(lastID), newUUID, nil
	}

	// Outro erro de DB
	return 0, "", err
}

func (s *Service) getOrCreateAuthor(tx *sqlx.Tx, name string) (int, string, error) {
	// Implementação similar ao getOrCreateCategory
	var id int
	var currentUUID string
	query := "SELECT id, uuid FROM authors WHERE name = ?"
	err := tx.QueryRowx(query, name).Scan(&id, &currentUUID)
	if err == nil {
		// Categoria encontrada: retornamos os dados dela.
		return id, currentUUID, nil
	}

	// Verificamos se o erro foi apenas que a linha não existe.
	if errors.Is(err, sql.ErrNoRows) {
		// Categoria não existe, vamos criar:

		newUUID := uuid.New().String() // Gerar NOVO UUID

		// 2. Inserir um novo registro, incluindo o UUID.
		result, err := tx.Exec("INSERT INTO authors (name, uuid) VALUES (?, ?)", name, newUUID)
		if err != nil {
			return 0, "", err
		}

		lastID, _ := result.LastInsertId()

		// 3. Retornar os IDs do novo registro.
		return int(lastID), newUUID, nil
	}

	// Outro erro de DB
	return 0, "", err

}

// Método para buscar um livro por ID
func (s *Service) GetBookByID(bookUUID string) (*Book, error) {
	var book Book

	query := `
        SELECT
            b.uuid, b.title, b.published_at, b.isbn, b.created_at,
            a.uuid AS "author.uuid", a.name AS "author.name",
            c.uuid AS "category.uuid", c.name AS "category.name"
        FROM books b
        JOIN authors a ON b.author_id = a.id
        JOIN categories c ON b.category_id = c.id
        WHERE b.uuid = ? -- Busca pelo UUID
    `
	// Passamos o bookUUID como parâmetro
	err := s.DB.Get(&book, query, bookUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &book, nil
}

// Método para buscar todos os livros
func (s *Service) GetAllBooks() ([]Book, error) {
	var books []Book

	query := `
        SELECT
            b.uuid, b.title, b.published_at, b.isbn, b.created_at,
            a.uuid AS "author.uuid", a.name AS "author.name",
            c.uuid AS "category.uuid", c.name AS "category.name"
        FROM books b
        JOIN authors a ON b.author_id = a.id
        JOIN categories c ON b.category_id = c.id
        ORDER BY b.created_at DESC
    `
	// Usamos Select do sqlx para buscar uma lista de structs
	err := s.DB.Select(&books, query)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *Service) DeleteBookByID(bookUUID string) error {
	// A exclusão não precisa de transação, pois é uma única operação.
	query := `DELETE FROM books WHERE uuid = ?`

	result, err := s.DB.Exec(query, bookUUID)
	if err != nil {
		return err // Retorna erro de DB
	}

	// Verifica quantas linhas foram afetadas para garantir que o livro existia
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound // Usa o erro definido anteriormente
	}

	return nil
}
