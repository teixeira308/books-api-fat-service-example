package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"example.com/books-api/internal/service"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input service.CreateBookInput

	// 1. Decodificar a requisição diretamente na struct de entrada do Serviço
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// 2. Chamar o método do Fat Service
	book, err := app.service.CreateBook(&input)

	// 3. Tratar os erros retornados pelo Serviço
	if err != nil {
		if errors.Is(err, service.ErrFailedValidation) {
			// Erro de validação: retorna 422 (Unprocessable Entity) com os erros
			app.failedValidation(w, r, input.ValidationErrors)
		} else {
			// Outro erro (DB, etc.): retorna 500 (Internal Server Error)
			app.serverError(w, r, err)
		}
		return
	}

	// 4. Sucesso: retorna 201 (Created) com o novo livro
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// -----------------------------------------------------
// Funções auxiliares (simplificadas)
// (Em um projeto real, você as criaria em um arquivo utilitário)
// -----------------------------------------------------

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("Bad Request: %v", err)
	http.Error(w, "Bad Request", http.StatusBadGateway)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Printf("Server Error: %v", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.serverError(w, nil, err)
	}
}

// Helper para ler o ID do path da URL (Ex: /books/123)
func (app *application) readIDParam(r *http.Request) (int, error) {
	// A rota é GET /books/{id}. O path contém "/books/123"
	// O mux do Go (a partir de 1.22) permite extrair variáveis de caminho.

	// No Go 1.22+ com http.ServeMux:
	idStr := r.PathValue("id")

	// Se você estiver usando uma versão anterior ou um router diferente,
	// você pode precisar usar strings.TrimPrefix para extrair a parte final do path:
	// path := r.URL.Path
	// parts := strings.Split(path, "/")
	// idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

// Handler para buscar um livro específico
func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	bookUUID := r.PathValue("id")

	if bookUUID == "" {
		app.badRequest(w, r, errors.New("book ID (UUID) must be provided"))
		return
	}

	// 2. Chamar o método do Fat Service
	book, err := app.service.GetBookByID(bookUUID)

	// 3. Tratar erros
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// 4. Sucesso: retorna 200 OK com o livro
	app.writeJSON(w, http.StatusOK, book)
}

// Handler para buscar todos os livros
func (app *application) getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Chamar o método do Fat Service
	books, err := app.service.GetAllBooks()

	// 2. Tratar erros
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// 3. Sucesso: retorna 200 OK com a lista (pode ser vazia)
	app.writeJSON(w, http.StatusOK, books)
}

// Handler para deletar um livro específico
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Ler o ID do parâmetro de caminho (reutiliza o helper)
	bookUUID := r.PathValue("id")
	if bookUUID == "" {
		app.badRequest(w, r, errors.New("book ID (UUID) must be provided"))
		return
	}

	// 2. Chamar o método do Fat Service
	err := app.service.DeleteBookByID(bookUUID)

	// 3. Tratar erros
	if err != nil {
		if errors.Is(err, service.ErrRecordNotFound) {
			http.NotFound(w, r) // 404 Not Found se o livro não existir
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// 4. Sucesso: retorna 204 No Content (padrão para exclusão bem-sucedida)
	w.WriteHeader(http.StatusNoContent)
}
