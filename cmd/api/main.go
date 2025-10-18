package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"example.com/books-api/internal/service"
	"github.com/jmoiron/sqlx"

	// Importa o driver SQLite
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger  *log.Logger
	service *service.Service
}

func main() {
	dsn := flag.String("dsn", "./books.sqlite", "SQLite DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// 1. Inicializar o Banco de Dados (SQLite)
	db, err := sqlx.Connect("sqlite3", *dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Configuração inicial do DB (Criar tabelas)
	if err := createSchema(db); err != nil {
		logger.Fatal(err)
	}

	// 2. Inicializar o Fat Service
	appService := &service.Service{DB: db}

	app := &application{
		logger:  logger,
		service: appService,
	}

	// 3. Configurar Rotas
	mux := http.NewServeMux()
	mux.HandleFunc("POST /books", app.createBookHandler)
	// mux.HandleFunc("GET /books/{id}", app.getBookHandler) // Exemplo de outras rotas

	// Rotas GET (Busca)
	// Busca um livro por ID. Note a sintaxe com a variável {id}
	mux.HandleFunc("GET /books/{id}", app.getBookHandler)

	// Busca todos os livros
	mux.HandleFunc("GET /books", app.getAllBooksHandler)

	// Rota DELETE (Deletar)
	mux.HandleFunc("DELETE /books/{id}", app.deleteBookHandler)

	// 4. Iniciar o Servidor
	logger.Print("starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	logger.Fatal(err)
}

// Função para criar as tabelas no DB SQLite
func createSchema(db *sqlx.DB) error {
	schema := `
        CREATE TABLE IF NOT EXISTS authors (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL, -- NOVO CAMPO UUID
            name TEXT NOT NULL UNIQUE
        );
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL, -- NOVO CAMPO UUID
            name TEXT NOT NULL UNIQUE
        );
        CREATE TABLE IF NOT EXISTS books (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid TEXT UNIQUE NOT NULL, -- NOVO CAMPO UUID
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
