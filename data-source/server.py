from flask import Flask, jsonify
import csv
import os

app = Flask(__name__)

CSV_FILE = '/app/orders.csv'

def read_orders():
    """Lê o arquivo CSV e retorna os dados como lista de dicionários"""
    orders = []
    
    if not os.path.exists(CSV_FILE):
        return orders
    
    with open(CSV_FILE, 'r', encoding='utf-8') as file:
        reader = csv.DictReader(file, delimiter=';')
        
        for row in reader:
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

@app.route('/')
def get_orders():
    """Endpoint GET que retorna todos os pedidos do CSV"""
    try:
        orders = read_orders()
        return jsonify(orders), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/health')
def health():
    """Endpoint de health check"""
    return {'status': 'healthy'}, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=3000, debug=True)


