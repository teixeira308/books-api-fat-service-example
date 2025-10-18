package service_test

import (
	"testing"

	"example.com/books-api/internal/service" // Seu pacote de serviço
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateBook_Success(t *testing.T) {
	// 1. SETUP: Conexão DB em memória
	db, err := sqlx.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close() // Fecha e descarta o DB após o teste

	if err := createSchemaForTest(db); err != nil {
		t.Fatalf("Falha ao criar o schema do DB em memória: %v", err)
	}

	// Cria a instância do Serviço
	svc := service.NewService(db)

	// 2. EXECUTAR O TESTE
	input := &service.CreateBookInput{
		Title:        "Livro Teste",
		AuthorName:   "Autor Teste",
		CategoryName: "Categoria Teste",
		PublishedAt:  20240101,
		ISBN:         "123-4567890123",
		// ... outros campos
	}

	book, err := svc.CreateBook(input)

	// 3. VERIFICAR
	if err != nil {
		t.Fatalf("Esperava sucesso, mas obteve erro: %v", err)
	}
	if book.Title != "Livro Teste" {
		t.Errorf("Título incorreto.")
	}
	// Você também verificaria se o book.UUID foi gerado corretamente.

	// 4. VERIFICAÇÃO FINAL NO DB (Opcional, mas comum)
	var count int
	db.Get(&count, "SELECT count(id) FROM books WHERE title = ?", "Livro Teste")
	if count != 1 {
		t.Errorf("Esperava 1 livro, encontrou %d", count)
	}
}

// Exemplo de código para seu arquivo service_test.go

// A função createSchema que você tem em cmd/api/main.go
func createSchemaForTest(db *sqlx.DB) error {
	schema := `
        CREATE TABLE IF NOT EXISTS authors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL,
            name TEXT NOT NULL UNIQUE
        );
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL,
            name TEXT NOT NULL UNIQUE
        );
        CREATE TABLE IF NOT EXISTS books (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL,
            title TEXT NOT NULL,
            author_id INTEGER,
            category_id INTEGER,
            published_at INTEGER,
            isbn TEXT,
            created_at DATETIME,
            FOREIGN KEY(author_id) REFERENCES authors(id),
            FOREIGN KEY(category_id) REFERENCES categories(id)
        );
    `
	_, err := db.Exec(schema)
	return err
}
