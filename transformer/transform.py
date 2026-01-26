import os
import psycopg2
from psycopg2.extras import RealDictCursor
from flask import Flask, jsonify
from flask_cors import CORS

def get_database_connection():
    """Conecta ao PostgreSQL usando DATABASE_URL"""
    database_url = os.getenv("DATABASE_URL") # lê a variável de ambiente DATABASE_URL
    if not database_url:
        raise ValueError("DATABASE_URL não configurada")
    
    # Adicionar sslmode=disable se não estiver presente (PostgreSQL local não usa SSL)
    if "sslmode" not in database_url:
        if "?" in database_url:
            database_url += "&sslmode=disable"
        else:
            database_url += "?sslmode=disable"
    
    # Conectar ao banco
    conn = psycopg2.connect(database_url) # conecta ao banco de dados usando a URL de conexão
    return conn

def setup_aggregated_schema(conn):
    """Cria o schema aggregated e a tabela daily_metrics se não existirem"""
    with conn.cursor() as cur: # cursor é um objeto que permite executar consultas SQL
        # Criar schema aggregated se não existir
        cur.execute("CREATE SCHEMA IF NOT EXISTS aggregated")
        
        # Criar tabela aggregated.daily_metrics se não existir
        create_table_sql = """
            CREATE TABLE IF NOT EXISTS aggregated.daily_metrics (
                id SERIAL PRIMARY KEY,
                date DATE NOT NULL,
                status VARCHAR(50) NOT NULL,
                payment_method VARCHAR(50) NOT NULL,
                total_orders INTEGER NOT NULL,
                total_value NUMERIC(10, 2) NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                UNIQUE(date, status, payment_method)  -- garante que não há duplicidade de dados na tabela aggregated.daily_metrics
            )
        """
        cur.execute(create_table_sql) # executa o SQL de criação da tabela
        
        conn.commit()
        print("✅ Schema aggregated e tabela daily_metrics verificados/criados")

def aggregate_data(conn):
    """Lê dados de raw_data.orders e agrega por data, status e payment_method"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur: # retorna linhas como dicionários, para o acesso ser dado por nome de coluna, em vez de índice
        # Query de agregação, query é uma consulta SQL que retorna os dados agregados por data, status e payment_method, é uma query de seleção (SELECT)
        aggregation_sql = """
            SELECT -- seleciona quais colunas serão retornadas
                DATE(created_at) as date,
                status,
                payment_method,
                COUNT(*) as total_orders,
                SUM(value) as total_value
            FROM raw_data.orders
            GROUP BY DATE(created_at), status, payment_method  -- agrupa os que tem o mesmo date, status e payment_method
            ORDER BY date, status, payment_method  -- ordena por date, status e payment_method
        """
        
        cur.execute(aggregation_sql) # executa o SQL de agregação
        aggregated_data = cur.fetchall() # retorna as linhas resultantes da execução do SQL de agregação
        
        print(f"✅ {len(aggregated_data)} grupos de dados agregados encontrados")
        return aggregated_data

def insert_aggregated_data(conn, aggregated_data): # recebe a conexão e os dados agregados e atualiza a tabela aggregated.daily_metrics
    """Insere os dados agregados na tabela aggregated.daily_metrics"""
    if not aggregated_data:
        print("⚠️  Nenhum dado para inserir")
        return 0
    
    with conn.cursor() as cur: # cursor é um objeto que permite executar consultas SQL
        # Preparar statement de inserção, insere os dados agregados na tabela aggregated.daily_metrics.
        insert_sql = """
            INSERT INTO aggregated.daily_metrics
                (date, status, payment_method, total_orders, total_value)
            VALUES (%s, %s, %s, %s, %s) -- placeholders para os valores a serem inseridos na tabela aggregated.daily_metrics
            ON CONFLICT (date, status, payment_method) -- se já existir uma linha com a mesma data, status e payment_method, atualiza os valores para os novos que está tentando inserir, ou mantém o mesmo se não for um dado novo (transformer roda novamente sem dados novos)
            DO UPDATE SET
                total_orders = EXCLUDED.total_orders, -- excluded é o valor que você quer inserir, com os mesmos valores de date, status e payment_method que já existem
                total_value = EXCLUDED.total_value,
                created_at = CURRENT_TIMESTAMP
        """
        
        inserted = 0
        for row in aggregated_data: # para cada linha na lista de dados agregados, insere na tabela aggregated.daily_metrics; Executa o statement preparado para cada linha
            try:
                cur.execute(
                    insert_sql, # insere os valores nos placeholders
                    (
                        row['date'],
                        row['status'],
                        row['payment_method'],
                        row['total_orders'],
                        float(row['total_value'])
                    )
                )
                inserted += 1
            except Exception as e: # se houver erro, imprime o erro e continua para a próxima linha
                print(f"⚠️  Erro ao inserir linha: {e}")
                continue
        
        conn.commit() # confirma a transação, ou seja, insere as linhas na tabela aggregated.daily_metrics. antes disso, ficam como pendentes
        return inserted # retorna o número de linhas inseridas

def run_transformation():
    """Executa a transformação de dados"""
    try:
        # Conectar ao PostgreSQL
        print("\n📡 Conectando ao PostgreSQL...")
        conn = get_database_connection()
        print("✅ Conectado ao PostgreSQL")
        
        # Configurar schema e tabela
        print("\n🏗️  Configurando schema aggregated...")
        setup_aggregated_schema(conn)
        
        # Agregar dados
        print("\n📊 Agregando dados de raw_data.orders...")
        aggregated_data = aggregate_data(conn)
        
        # Inserir dados agregados
        print("\n💾 Inserindo dados agregados em aggregated.daily_metrics...")
        inserted = insert_aggregated_data(conn, aggregated_data)
        print(f"✅ {inserted} registros inseridos/atualizados com sucesso")
        
        # Fechar conexão com o banco de dados
        conn.close()
        
        print("\n=== Transformação concluída com sucesso ===")
        return inserted # retorna o número de linhas inseridas
        
    except Exception as e:
        print(f"\n❌ Erro: {e}")
        raise

# Criar aplicação Flask
app = Flask(__name__) # flask é um framework da API para Python
CORS(app)  # Habilitar CORS

@app.route('/') # rota get para a raiz do serviço
def hello():
    return {
        'service': 'transformer',
        'status': 'running',
        'message': 'Serviço de Transformação de Dados'
    }

@app.route('/health') # rota get para o health check
def health():
    return {'status': 'healthy'}, 200

@app.route('/transform', methods=['POST']) #rota post para executar a transformação
def transform():
    """Endpoint HTTP para executar a transformação"""
    try:
        print("\n=== Transformação disparada via HTTP ===")
        inserted = run_transformation() # executa a transformação e retorna o número de linhas inseridas
        return jsonify({
            'success': True,
            'message': 'Transformação executada com sucesso',
            'inserted': inserted
        }), 200
    except Exception as e: # se houver erro, retorna o erro
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500

def main():
    """Função principal - executa transformação uma vez na inicialização"""
    print("=== Serviço de Transformação de Dados iniciado ===")
    try:
        run_transformation()
    except Exception as e:
        print(f"\n❌ Erro: {e}")
        raise

if __name__ == '__main__':
    # Se executado diretamente (não via import), iniciar servidor HTTP
    port = int(os.getenv('PORT', '8080'))
    print(f"\n🚀 Servidor HTTP iniciado na porta {port}")
    print("Endpoints disponíveis:")
    print("  - GET  /health    - Health check")
    print("  - POST /transform  - Executar transformação")
    app.run(host='0.0.0.0', port=port, debug=False)
