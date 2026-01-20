package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Order representa a estrutura de um pedido recebido da API
type Order struct {
	OrderID       string  `json:"order_id"`
	CreatedAt     string  `json:"created_at"`
	Status        string  `json:"status"`
	Value         float64 `json:"value"`
	PaymentMethod string  `json:"payment_method"`
}

func main() {
	fmt.Println("=== Pipeline de Dados iniciado ===")

	// Obter URLs das variáveis de ambiente
	dataSourceURL := os.Getenv("DATA_SOURCE_URL")
	if dataSourceURL == "" {
		dataSourceURL = "http://data-source:3000"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL não configurada")
	}

	// Adicionar sslmode=disable se não estiver presente (PostgreSQL local não usa SSL)
	if !strings.Contains(databaseURL, "sslmode") {
		if strings.Contains(databaseURL, "?") {
			databaseURL += "&sslmode=disable"
		} else {
			databaseURL += "?sslmode=disable"
		}
	}

	fmt.Printf("Data Source URL: %s\n", dataSourceURL)
	fmt.Printf("Database URL: %s\n", databaseURL)

	// Conectar ao PostgreSQL
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}
	defer db.Close()

	// Testar conexão
	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao fazer ping no PostgreSQL: %v", err)
	}
	fmt.Println("✅ Conectado ao PostgreSQL")

	// Criar schema e tabela se não existirem
	if err := setupDatabase(db); err != nil {
		log.Fatalf("Erro ao configurar banco de dados: %v", err)
	}
	fmt.Println("✅ Schema e tabela verificados/criados")

	// Buscar dados do Data Source
	fmt.Println("\n📥 Buscando dados do Data Source...")
	orders, err := fetchOrders(dataSourceURL)
	if err != nil {
		log.Fatalf("Erro ao buscar pedidos: %v", err)
	}
	fmt.Printf("✅ %d pedidos recebidos do Data Source\n", len(orders))

	// Inserir dados no banco
	fmt.Println("\n💾 Inserindo dados no PostgreSQL...")
	inserted, err := insertOrders(db, orders)
	if err != nil {
		log.Fatalf("Erro ao inserir pedidos: %v", err)
	}
	fmt.Printf("✅ %d pedidos inseridos com sucesso\n", inserted)

	fmt.Println("\n=== Pipeline concluído com sucesso ===")
}

// setupDatabase cria o schema raw_data e a tabela orders se não existirem
func setupDatabase(db *sql.DB) error {
	// Criar schema raw_data se não existir
	_, err := db.Exec("CREATE SCHEMA IF NOT EXISTS raw_data")
	if err != nil {
		return fmt.Errorf("erro ao criar schema: %w", err)
	}

	// Criar tabela raw_data.orders se não existir
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS raw_data.orders (
			id SERIAL PRIMARY KEY,
			order_id VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL,
			status VARCHAR(50) NOT NULL,
			value NUMERIC(10, 2) NOT NULL,
			payment_method VARCHAR(50) NOT NULL,
			created_at_pipeline TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela: %w", err)
	}

	return nil
}

// fetchOrders busca os pedidos da API do Data Source
func fetchOrders(url string) ([]Order, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code não OK: %d", resp.StatusCode)
	}

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return orders, nil
}

// insertOrders insere os pedidos no banco de dados
func insertOrders(db *sql.DB, orders []Order) (int, error) {
	if len(orders) == 0 {
		return 0, nil
	}

	// Preparar statement SQL para inserção
	stmt, err := db.Prepare(`
		INSERT INTO raw_data.orders (order_id, created_at, status, value, payment_method)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (order_id) DO NOTHING
	`)
	if err != nil {
		return 0, fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	inserted := 0
	for _, order := range orders {
		// Converter created_at de string para time.Time
		createdAt, err := time.Parse(time.RFC3339, order.CreatedAt)
		if err != nil {
			log.Printf("⚠️  Erro ao parsear created_at '%s': %v", order.CreatedAt, err)
			continue
		}

		// Inserir no banco
		result, err := stmt.Exec(
			order.OrderID,
			createdAt,
			order.Status,
			order.Value,
			order.PaymentMethod,
		)
		if err != nil {
			log.Printf("⚠️  Erro ao inserir pedido %s: %v", order.OrderID, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			inserted++
		}
	}

	return inserted, nil
}
