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

// PipelineResponse representa a resposta do endpoint /trigger
type PipelineResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Inserted  int    `json:"inserted"`
	Total     int    `json:"total"`
	Timestamp string `json:"timestamp"`
}

var db *sql.DB
var dataSourceURL string  // var global
var transformerURL string // var global

func main() {
	fmt.Println("=== Pipeline de Dados iniciado ===")

	// Obter URLs das variáveis de ambiente, "os" verifica se existe antes de ler
	dataSourceURL = os.Getenv("DATA_SOURCE_URL")
	if dataSourceURL == "" {
		dataSourceURL = "http://data-source:3000"
	}

	transformerURL = os.Getenv("TRANSFORMER_URL")
	if transformerURL == "" {
		transformerURL = "http://transformer:8080/transform"
	}

	databaseURL := os.Getenv("DATABASE_URL") // var local
	if databaseURL == "" {
		log.Fatal("DATABASE_URL não configurada") // encerra o programa
	}

	// desabilitar ssl, pois a conexão com o PostgreSQL local não usa SSL (comunicação não atravessa internet)
	if !strings.Contains(databaseURL, "sslmode") {
		if strings.Contains(databaseURL, "?") { // ? para primeiro parâmetro, & parâmetros adicionais
			databaseURL += "&sslmode=disable"
		} else {
			databaseURL += "?sslmode=disable"
		}
	}

	fmt.Printf("Data Source URL: %s\n", dataSourceURL)
	fmt.Printf("Transformer URL: %s\n", transformerURL)
	fmt.Printf("Database URL: %s\n", databaseURL)

	// Conectar ao PostgreSQL
	var err error
	db, err = sql.Open("postgres", databaseURL) // sql.Open é uma função que abre uma conexão com o PostgreSQL (sem API, conexão direta via driver de banco de dados)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}
	defer db.Close()

	// Testar conexão
	if err := db.Ping(); err != nil { // ping != 0 significa erro
		log.Fatalf("Erro ao fazer ping no PostgreSQL: %v", err)
	}
	fmt.Println("✅ Conectado ao PostgreSQL")

	// Criar schema e tabela se não existirem
	if err := setupDatabase(db); err != nil { // setupDatabase é uma função que cria o schema e a tabela se não existirem
		log.Fatalf("Erro ao configurar banco de dados: %v", err)
	}
	fmt.Println("✅ Schema e tabela verificados/criados")

	// Configurar rotas HTTP
	// handler é uma função que processa a requisição e escreve a resposta
	http.HandleFunc("/health", healthHandler)   // registra handler para GET /health
	http.HandleFunc("/trigger", triggerHandler) // registra handler para POST /trigger

	// Iniciar servidor HTTP
	port := os.Getenv("PORT") // port é a porta do servidor HTTP
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 Servidor HTTP iniciado na porta %s\n", port)
	fmt.Println("Endpoints disponíveis:")
	fmt.Println("  - GET  /health  - Health check")
	fmt.Println("  - POST /trigger - Disparar ingestão de dados")

	log.Fatal(http.ListenAndServe(":"+port, nil)) // inicia o servidor na porta ou encerra o programa se houver erro
}

func healthHandler(w http.ResponseWriter, r *http.Request) { // w (response writer) é o objeto que escreve a resposta, r (request) é o objeto que representa a requisição
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // se o método de requisição não é GET, retorna erro
		return
	}

	w.Header().Set("Content-Type", "application/json")                // define o header como application/json, informando que o conteúdo é JSON
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}) // cria um encoder json que converte objeto Go em JSON e escreve em "w" a resposta "{"status": "healthy"}"
}

func triggerHandler(w http.ResponseWriter, r *http.Request) { // função que dispara a execução do pipeline e escreve a resposta
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // se o método de requisição não é POST, retorna erro
		return
	}

	fmt.Println("\n=== Pipeline disparado via HTTP ===")

	// Executar pipeline
	inserted, total, err := runPipeline() // executa o processo de ingestão de dados no PostgreSQL

	response := PipelineResponse{ //
		Success:   err == nil, // sucess = true se err é nil
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err != nil {
		response.Message = fmt.Sprintf("Erro ao executar pipeline: %v", err)
		w.Header().Set("Content-Type", "application/json") // informa que a resposta será JSIN, definindo o header como application/json
		w.WriteHeader(http.StatusInternalServerError)      // escreve o status code 500 (Internal Server Error)
		json.NewEncoder(w).Encode(response)                // cria um encoder json que converte objeto "response" de Go para JSON e escreve em "w" a resposta
		return
	}

	response.Message = "Pipeline executado com sucesso"
	response.Inserted = inserted
	response.Total = total

	w.Header().Set("Content-Type", "application/json") // informa que a resposta será JSON
	json.NewEncoder(w).Encode(response)                // converte objeto para JSON
}

func runPipeline() (int, int, error) { // retorna 2 int e 1 error
	// Buscar dados do Data Source
	fmt.Println("\n📥 Buscando dados do Data Source...")
	orders, err := fetchOrders(dataSourceURL) // fetchOrders é uma função que busca os pedidos da API do Data Source
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao buscar pedidos: %w", err)
	}
	fmt.Printf("✅ %d pedidos recebidos do Data Source\n", len(orders)) // qtd de pedidos recebidos

	// Inserir dados no banco
	fmt.Println("\n💾 Inserindo dados no PostgreSQL...")
	inserted, err := insertOrders(db, orders) // insertOrders é uma função que insere os pedidos no banco de dados
	if err != nil {
		return 0, len(orders), fmt.Errorf("erro ao inserir pedidos: %w", err)
	}
	fmt.Printf("✅ %d pedidos inseridos com sucesso\n", inserted) // qtd de pedidos inseridos

	// Chamar transformer para agregar dados
	if inserted > 0 {
		fmt.Println("\n🔄 Chamando transformer para agregar dados...")
		if err := callTransformer(transformerURL); err != nil { // callTransformer é uma função que chama o serviço transformer via HTTP
			log.Printf("⚠️  Erro ao chamar transformer: %v", err)
			// Não falhar o pipeline se o transformer falhar
		} else {
			fmt.Println("✅ Transformer executado com sucesso")
		}
	}

	fmt.Println("\n=== Pipeline concluído com sucesso ===")
	return inserted, len(orders), nil
}

// setupDatabase cria o schema raw_data e a tabela orders se não existirem
func setupDatabase(db *sql.DB) error {
	// Criar schema raw_data se não existir
	_, err := db.Exec("CREATE SCHEMA IF NOT EXISTS raw_data") // db.Exec é uma função que executa uma consulta SQL
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

	_, err = db.Exec(createTableSQL) // executa o SQL de criação da tabela
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

	resp, err := client.Get(url) // faz uma requisição GET para a URL para obter os dados do Data Source
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // status ok = 200
		return nil, fmt.Errorf("status code não OK: %d", resp.StatusCode)
	}

	var orders []Order                                                 // orders é um slice de Order. Order é uma ficha de pedido, e []Order é uma pasta com várias fichas
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil { // decodifica o JSON da resposta, passando de JSON para Go
		return nil, fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	return orders, nil // retorna os pedidos em formato Go ou erro se houver
}

// insertOrders insere os pedidos no banco de dados
func insertOrders(db *sql.DB, orders []Order) (int, error) {
	if len(orders) == 0 {
		return 0, nil
	}

	// Preparar statement (stmt) SQL para inserção, cria um template SQL que será executado posteriormente com os valores passados
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
	for _, order := range orders { // para cada pedido, executa o statement preparado
		// Converter created_at de string para time.Time
		createdAt, err := time.Parse(time.RFC3339, order.CreatedAt) // converte a string para time.Time
		if err != nil {
			log.Printf("⚠️  Erro ao parsear created_at '%s': %v", order.CreatedAt, err) // parsear é transformar texto bruto em dado estruturado
			continue
		}

		// Inserir no banco
		result, err := stmt.Exec( // executa o statement preparado, ou seja, preenche os valores do template SQL com os valores do pedido
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

		rowsAffected, _ := result.RowsAffected() // obtém o número de linhas afetadas pela execução do statement
		if rowsAffected > 0 {
			inserted++
		}
	}

	return inserted, nil
}

// callTransformer chama o serviço transformer via HTTP
func callTransformer(url string) error {
	client := &http.Client{ // acessa o endpoint do transformer via HTTP
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(url, "application/json", nil) // faz uma requisição POST (pois executa transformação nos dados) para a URL
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code não OK: %d", resp.StatusCode)
	}

	return nil
}
