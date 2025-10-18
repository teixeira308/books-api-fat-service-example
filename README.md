# 📚 Books API - Go Fat Service Architecture

Uma API REST simples e robusta desenvolvida em Go, utilizando SQLite como banco de dados. Este projeto serve como um estudo prático e uma demonstração da arquitetura **"Fat Service"**, priorizando a simplicidade, a coesão do código e a filosofia "less is more" do Go.

## 🚀 Arquitetura Utilizada: Fat Service

Em contraste com a arquitetura tradicional Handler-Service-Repository, esta API adota o padrão **Fat Service** (apelido para Service Object + Persistência direta) para reduzir *boilerplate* e simplificar o gerenciamento de transações.

| Camada | Localização | Responsabilidade Principal | 
 | ----- | ----- | ----- | 
| **Aplicação** | `cmd/api/handlers.go` | Recebe requisições HTTP e formata respostas. | 
| **Fat Service** | `internal/service/*.go` | **Lógica de Negócios** (Validação, Geração de UUIDs) e **Persistência de Dados** (SQL Direto). | 

**Destaque Arquitetural:** O Service concentra todas as responsabilidades do domínio, eliminando a Camada de Repositório e tornando a gestão de transações (`sql.Tx`) mais direta e legível, alinhada à filosofia de simplicidade de Go.

## ✨ Características Chave

* **Identificadores Seguros (UUIDs):** Utiliza UUIDs (`string`) como identificadores públicos da API em todas as rotas (`GET`, `DELETE`), prevenindo a enumeração de recursos.

* **Banco de Dados Embarcado:** Utiliza **SQLite** (`books.sqlite`) com `sqlx`.

* **Testes de Integração Rápidos:** O Service é testado contra um DB SQLite em modo **`file::memory:`** nos testes, garantindo a confiança no SQL e no mapeamento do `sqlx` (o *trade-off* do Fat Service).

## ⚙️ Configuração e Execução

### Pré-requisitos

* Go (Versão 1.20+)

### 1. Clonar o Repositório

git clone [URL_DO_SEU_REPOSITORIO] cd books-api


### 2. Instalar Dependências

Certifique-se de que as dependências, como `sqlx` e `go-sqlite3`, estão instaladas:

go mod tidy


### 3. Rodar a API

A API é compilada e executada a partir do diretório `cmd/api`. A primeira execução criará o arquivo `books.sqlite`.

go run ./cmd/api


A API estará disponível em `http://localhost:4000`.

## 🧪 Testes

Execute todos os testes de integração e unitários do projeto:

go test ./...


## 🧭 Endpoints da API

Todos os IDs nas requisições e respostas são **UUIDs**.

| Método | Caminho | Descrição | 
 | ----- | ----- | ----- | 
| **POST** | `/books` | Cria um novo livro (e cria Autor/Categoria se necessário). | 
| **GET** | `/books` | Lista todos os livros. | 
| **GET** | `/books/{id}` | Busca um livro específico pelo UUID. | 
| **DELETE** | `/books/{id}` | Deleta um livro específico pelo UUID. | 

### Exemplo de Uso (Criação)

curl -X POST http://localhost:4000/books -H "Content-Type: application/json" -d ' { "title": "A Máquina do Tempo", "author_name": "H.G. Wells", "category_name": "Ficção Científica", "published_at": 1895, "isbn": "978-8575034637" } '
