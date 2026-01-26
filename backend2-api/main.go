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

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

// define estruturas de dados para as respostas das APIs

type MetricsResponse struct {
	Filters            Filters            `json:"filters"`
	FinancialMetrics   FinancialMetrics   `json:"financial_metrics"`
	OperationalMetrics OperationalMetrics `json:"operational_metrics"`
}

type Filters struct {
	StartDate     string `json:"start_date,omitempty"`
	EndDate       string `json:"end_date,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
}

type FinancialMetrics struct {
	ApprovedRevenue  float64 `json:"approved_revenue"`
	PendingRevenue   float64 `json:"pending_revenue"`
	CancelledRevenue float64 `json:"cancelled_revenue"`
}

type OperationalMetrics struct { // estrutura métricas operacionais
	ApprovedOrders  int `json:"approved_orders"`
	PendingOrders   int `json:"pending_orders"`
	CancelledOrders int `json:"cancelled_orders"`
}

type TimeSeriesResponse struct { // estrutura série temporal
	Filters Filters           `json:"filters"`
	Data    []TimeSeriesPoint `json:"data"`
}

type TimeSeriesPoint struct { // estrutura para um ponto (dia) da série temporal
	Date             string  `json:"date"`
	ApprovedRevenue  float64 `json:"approved_revenue"`
	PendingRevenue   float64 `json:"pending_revenue"`
	CancelledRevenue float64 `json:"cancelled_revenue"`
	ApprovedOrders   int     `json:"approved_orders"`
	PendingOrders    int     `json:"pending_orders"`
	CancelledOrders  int     `json:"cancelled_orders"`
}

var jwtSecret string // variável global para a chave JWT

func main() {
	// Obter JWT_SECRET da variável de ambiente (deve ser a mesma do backend1-auth)
	jwtSecret = os.Getenv("JWT_SECRET") // pega a chave JWT da variável de ambiente (configuração do sistema, não do código)
	if jwtSecret == "" {                // se estiver vazia, usa o valor padrão
		jwtSecret = "minha-chave-secreta-jwt-super-segura" // Fallback para desenvolvimento
		log.Println("⚠️  JWT_SECRET não configurada, usando valor padrão")
	}

	// Configurar rotas para expor endpoints
	http.HandleFunc("/", corsMiddleware(helloHandler)) // todas são protegidas pelo CORS
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/api/metrics", corsMiddleware(verifyTokenMiddleware(metricsHandler)))                // métricas são protegidas pelo JWT
	http.HandleFunc("/api/metrics/time-series", corsMiddleware(verifyTokenMiddleware(timeSeriesHandler))) // séries temporais também

	fmt.Println("Backend 2 API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// corsMiddleware adiciona headers CORS às respostas
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc { // middleware (controle de acesso, faz verificações) para adicionar headers CORS às respostas
	return func(w http.ResponseWriter, r *http.Request) {
		// Permitir origem do frontend
		w.Header().Set("Access-Control-Allow-Origin", "*")                            // permite acesso de qualquer origem; qualquer site pode chamar essa API
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // métodos permitidos
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // headers permitidos

		// Responder a requisições OPTIONS (preflight)
		if r.Method == "OPTIONS" { // options é um método que é usado para verificar se o servidor suporta o método de requisição
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// define o middleware que verifica se o token JWT é válido antes de permitir acesso
func verifyTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obter token do header Authorization
		authHeader := r.Header.Get("Authorization") // pega o token do header Authorization
		if authHeader == "" {                       // se o token não foi fornecido
			http.Error(w, `{"error": "Token não fornecido"}`, http.StatusUnauthorized)
			return
		}

		// Formato esperado: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" { // se o token não estiver no formato esperado
			http.Error(w, `{"error": "Formato de token inválido. Use: Bearer <token>"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1] // pega apenas o token do header Authorization

		// Verificar e decodificar o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // decodifica o token para obter os dados do usuário e validade do token
			// Verificar algoritmo de assinatura do token
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // se o algoritmo de assinatura não for HMAC
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"]) // retorna erro
			}
			return []byte(jwtSecret), nil // retorna a chave JWT
		})

		if err != nil { // se houver erro ao decodificar o token
			http.Error(w, fmt.Sprintf(`{"error": "Token inválido: %v"}`, err), http.StatusUnauthorized)
			return
		}

		if !token.Valid { // se o token não for válido
			http.Error(w, `{"error": "Token inválido"}`, http.StatusUnauthorized)
			return
		}

		// Token válido, continuar com a requisição
		next(w, r)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) { // define o handler para a rota raiz, endpoint retorna informações básicas do serviço
	response := map[string]string{ // cria mapa (dicionário) para a resposta
		"service": "backend2-api",
		"status":  "running",
		"message": "Hello from Backend 2 (Query API)",
	}

	w.Header().Set("Content-Type", "application/json") // define o header content-type como json
	json.NewEncoder(w).Encode(response)                // codifica o mapa em json e escreve na resposta
}

func healthHandler(w http.ResponseWriter, r *http.Request) { // define o handler para a rota /health, endpoint verifica a saúde do serviço
	w.Header().Set("Content-Type", "application/json")                // define o header content-type como json
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}) // cria um mapa com status e codifica o json
}

func getDB() (*sql.DB, error) { // função que abre e retorna uma conexão com o PostgreSQL ou erro
	databaseURL := os.Getenv("DATABASE_URL") // lê DATABASE_URL
	if databaseURL == "" {                   // se DATABASE_URL não estiver configurada (vazia)
		return nil, fmt.Errorf("DATABASE_URL não configurada")
	}

	// Adicionar sslmode=disable se não estiver presente
	if !strings.Contains(databaseURL, "sslmode") {
		if strings.Contains(databaseURL, "?") {
			databaseURL += "&sslmode=disable"
		} else {
			databaseURL += "?sslmode=disable"
		}
	}

	db, err := sql.Open("postgres", databaseURL) // abre uma conexão com o PostgreSQL
	if err != nil {                              // se houver erro ao abrir a conexão
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// metricsHandler retorna métricas agregadas (valores totais)
func metricsHandler(w http.ResponseWriter, r *http.Request) { // define o handler para a rota /api/metrics, endpoint retorna métricas agregadas (valores totais)
	if r.Method != http.MethodGet { // se o método não for GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obter parâmetros de query
	startDate := r.URL.Query().Get("start_date")         // pega o valor do parâmetro start_date da URL
	endDate := r.URL.Query().Get("end_date")             // pega o valor do parâmetro end_date
	paymentMethod := r.URL.Query().Get("payment_method") // pega o valor do parâmetro payment_method

	// Conectar ao banco
	db, err := getDB() // abre uma conexão com o PostgreSQL
	if err != nil {    // se houver erro ao abrir a conexão
		http.Error(w, fmt.Sprintf("Erro ao conectar ao banco: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Construir query base. query é uma instrução enviada ao banco de dados.
	query := `
		SELECT 
			status,
			SUM(total_orders) as total_orders, -- soma o total de pedidos
			SUM(total_value) as total_value -- soma o total de valor
		FROM aggregated.daily_metrics
		WHERE 1=1
	`

	args := []interface{}{} // slice vazio (pois não há argumentos na query base, então não sabemos quais filtros serão usados) para argumentos da query
	argIndex := 1

	// Adicionar filtros
	if startDate != "" { // se startDate não estiver vazio
		query += fmt.Sprintf(" AND date >= $%d", argIndex) // adiciona o filtro de data inicial à query
		args = append(args, startDate)                     // adiciona o valor de startDate ao slice de argumentos
		argIndex++
	}

	if endDate != "" {
		query += fmt.Sprintf(" AND date <= $%d", argIndex) // adiciona o filtro de data final à query
		args = append(args, endDate)                       // adiciona o valor de endDate ao slice de argumentos
		argIndex++
	}

	if paymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argIndex) // adiciona o filtro de método de pagamento à query
		args = append(args, paymentMethod)                          // adiciona o valor de paymentMethod ao slice de argumentos
		argIndex++
	}

	query += " GROUP BY status" // agrupa os resultados por status

	// Executar query
	rows, err := db.Query(query, args...) // executa a query
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao executar query: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Inicializar métricas
	metrics := MetricsResponse{ // cria uma estrutura para a resposta com os filtros e métricas
		Filters: Filters{
			StartDate:     startDate,
			EndDate:       endDate,
			PaymentMethod: paymentMethod,
		},
		FinancialMetrics:   FinancialMetrics{},
		OperationalMetrics: OperationalMetrics{},
	}

	// Processar resultados
	for rows.Next() {
		var status string
		var totalOrders int
		var totalValue float64

		if err := rows.Scan(&status, &totalOrders, &totalValue); err != nil { // se houver erro ao ler os resultados
			http.Error(w, fmt.Sprintf("Erro ao ler resultado: %v", err), http.StatusInternalServerError)
			return
		}

		switch status { // atribui os valores das métricas a cada status
		case "approved":
			metrics.FinancialMetrics.ApprovedRevenue = totalValue
			metrics.OperationalMetrics.ApprovedOrders = totalOrders
		case "pending":
			metrics.FinancialMetrics.PendingRevenue = totalValue
			metrics.OperationalMetrics.PendingOrders = totalOrders
		case "cancelled":
			metrics.FinancialMetrics.CancelledRevenue = totalValue
			metrics.OperationalMetrics.CancelledOrders = totalOrders
		}
	}

	w.Header().Set("Content-Type", "application/json") // define o header content-type como json
	json.NewEncoder(w).Encode(metrics)                 // codifica as métricas em json e escreve na resposta, que é enviada para o frontend
}

// timeSeriesHandler retorna séries temporais para gráficos
func timeSeriesHandler(w http.ResponseWriter, r *http.Request) { // função que define o handler para a rota /api/metrics/time-series, endpoint retorna séries temporais para gráficos
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Obter parâmetros de query
	startDate := r.URL.Query().Get("start_date") // pega o valor do parâmetro start_date
	endDate := r.URL.Query().Get("end_date")
	paymentMethod := r.URL.Query().Get("payment_method")

	// Conectar ao banco
	db, err := getDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao conectar ao banco: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Construir query para séries temporais
	query := `
		SELECT 
			date,
			SUM(CASE WHEN status = 'approved' THEN total_value ELSE 0 END) as approved_revenue, -- soma o total de valor para os pedidos aprovados, se não for aprovado é 0
			SUM(CASE WHEN status = 'pending' THEN total_value ELSE 0 END) as pending_revenue,
			SUM(CASE WHEN status = 'cancelled' THEN total_value ELSE 0 END) as cancelled_revenue,
			SUM(CASE WHEN status = 'approved' THEN total_orders ELSE 0 END) as approved_orders,
			SUM(CASE WHEN status = 'pending' THEN total_orders ELSE 0 END) as pending_orders,
			SUM(CASE WHEN status = 'cancelled' THEN total_orders ELSE 0 END) as cancelled_orders
		FROM aggregated.daily_metrics
		WHERE 1=1
	`

	args := []interface{}{} // slice vazio para argumentos da query
	argIndex := 1

	// Adicionar filtros
	if startDate != "" { // se startDate não estiver vazio
		query += fmt.Sprintf(" AND date >= $%d", argIndex) // adiciona o filtro de data inicial à query
		args = append(args, startDate)                     // adiciona o valor de startDate ao slice de argumentos
		argIndex++
	}

	if endDate != "" {
		query += fmt.Sprintf(" AND date <= $%d", argIndex) // adiciona o filtro de data final à query
		args = append(args, endDate)                       // adiciona o valor de endDate ao slice de argumentos
		argIndex++
	}

	if paymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argIndex) // adiciona o filtro de método de pagamento à query
		args = append(args, paymentMethod)                          // adiciona o valor de paymentMethod ao slice de argumentos
		argIndex++
	}

	query += " GROUP BY date ORDER BY date" // agrupa os resultados por data e ordena por data

	// Executar query
	rows, err := db.Query(query, args...) // executa a query com os argumentos
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao executar query: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Processar resultados
	var timeSeries []TimeSeriesPoint // slice vazio para os resultados da série temporal
	for rows.Next() {
		var point TimeSeriesPoint // variável para armazenar os pontos (dias) da série temporal
		var date time.Time        // cria uma variável para a data

		err := rows.Scan( // lê os resultados da query
			&date,                  // data
			&point.ApprovedRevenue, // lê a receita aprovada e armazena no ponto
			&point.PendingRevenue,
			&point.CancelledRevenue,
			&point.ApprovedOrders,
			&point.PendingOrders,
			&point.CancelledOrders,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao ler resultado: %v", err), http.StatusInternalServerError)
			return
		}

		point.Date = date.Format("2006-01-02") // formata a data no formato YYYY-MM-DD
		timeSeries = append(timeSeries, point) // adiciona o ponto à série temporal
	}

	// Criar resposta com filtros
	response := TimeSeriesResponse{ // cria uma estrutura para a resposta com os filtros e os pontos da série temporal
		Filters: Filters{
			StartDate:     startDate,     // data inicial
			EndDate:       endDate,       // data final
			PaymentMethod: paymentMethod, // método de pagamento
		},
		Data: timeSeries, // pontos da série temporal
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response) // codifica o response em json e escreve na resposta
}
