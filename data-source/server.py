from flask import Flask, jsonify
import csv
import os

app = Flask(__name__) # Inicializa o Flask

CSV_FILE = '/app/orders.csv' # Caminho do arquivo CSV

def read_orders(): # Função que lê o arquivo CSV e retorna os dados como lista de dicionários
    """Lê o arquivo CSV e retorna os dados como lista de dicionários"""
    orders = []
    
    if not os.path.exists(CSV_FILE): # se o arquivo não existe, retorna uma lista vazia
        return orders
    
    with open(CSV_FILE, 'r', encoding='utf-8') as file:
        reader = csv.DictReader(file, delimiter=';')
        
        for row in reader: # para cada linha do arquivo CSV, cria um dicionário com os dados da linha
            # Normaliza os dados conforme necessário
            order = {
                'order_id': row['order_id'],
                'created_at': row['created_at'],
                'status': row['status'],
                'value': float(row['value'].replace(',', '.')),  # Converte vírgula para ponto
                'payment_method': row['payment_method']
            }
            orders.append(order)
    
    return orders

@app.route('/') # Endpoint GET que retorna todos os pedidos do CSV
def get_orders(): # Função que retorna todos os pedidos do CSV
    """Endpoint GET que retorna todos os pedidos do CSV"""
    try:
        orders = read_orders() # chama a função read_orders para ler o arquivo CSV
        return jsonify(orders), 200 # retorna os dados como JSON
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/health')
def health():
    """Endpoint de health check"""
    return {'status': 'healthy'}, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000, debug=True)


