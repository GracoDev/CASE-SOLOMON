import os
import psycopg2
from psycopg2.extras import RealDictCursor

def get_database_connection():
    """Conecta ao PostgreSQL usando DATABASE_URL"""
    database_url = os.getenv("DATABASE_URL")
    if not database_url:
        raise ValueError("DATABASE_URL não configurada")
    
    # Adicionar sslmode=disable se não estiver presente (PostgreSQL local não usa SSL)
    if "sslmode" not in database_url:
        if "?" in database_url:
            database_url += "&sslmode=disable"
        else:
            database_url += "?sslmode=disable"
    
    # Conectar ao banco
    conn = psycopg2.connect(database_url)
    return conn

def setup_aggregated_schema(conn):
    """Cria o schema aggregated e a tabela daily_metrics se não existirem"""
    with conn.cursor() as cur:
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
                UNIQUE(date, status, payment_method)
            )
        """
        cur.execute(create_table_sql)
        
        conn.commit()
        print("✅ Schema aggregated e tabela daily_metrics verificados/criados")

def aggregate_data(conn):
    """Lê dados de raw_data.orders e agrega por data, status e payment_method"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        # Query de agregação
        aggregation_sql = """
            SELECT 
                DATE(created_at) as date,
                status,
                payment_method,
                COUNT(*) as total_orders,
                SUM(value) as total_value
            FROM raw_data.orders
            GROUP BY DATE(created_at), status, payment_method
            ORDER BY date, status, payment_method
        """
        
        cur.execute(aggregation_sql)
        aggregated_data = cur.fetchall()
        
        print(f"✅ {len(aggregated_data)} grupos de dados agregados encontrados")
        return aggregated_data

def insert_aggregated_data(conn, aggregated_data):
    """Insere os dados agregados na tabela aggregated.daily_metrics"""
    if not aggregated_data:
        print("⚠️  Nenhum dado para inserir")
        return 0
    
    with conn.cursor() as cur:
        # Preparar statement de inserção
        insert_sql = """
            INSERT INTO aggregated.daily_metrics 
                (date, status, payment_method, total_orders, total_value)
            VALUES (%s, %s, %s, %s, %s)
            ON CONFLICT (date, status, payment_method) 
            DO UPDATE SET
                total_orders = EXCLUDED.total_orders,
                total_value = EXCLUDED.total_value,
                created_at = CURRENT_TIMESTAMP
        """
        
        inserted = 0
        for row in aggregated_data:
            try:
                cur.execute(
                    insert_sql,
                    (
                        row['date'],
                        row['status'],
                        row['payment_method'],
                        row['total_orders'],
                        float(row['total_value'])
                    )
                )
                inserted += 1
            except Exception as e:
                print(f"⚠️  Erro ao inserir linha: {e}")
                continue
        
        conn.commit()
        return inserted

def main():
    print("=== Serviço de Transformação de Dados iniciado ===")
    
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
        
        # Fechar conexão
        conn.close()
        
        print("\n=== Transformação concluída com sucesso ===")
        
    except Exception as e:
        print(f"\n❌ Erro: {e}")
        raise

if __name__ == '__main__':
    main()
