from flask import Flask, request, jsonify
from flask_cors import CORS
from functools import wraps
import jwt
import datetime
import os
import requests

app = Flask(__name__)
CORS(app)  # Habilitar CORS para todas as rotas, permitindo que o frontend acesse o backend 1 sem problemas de CORS

# Configuração JWT
JWT_SECRET = os.getenv('JWT_SECRET', 'sua-chave-secreta-super-segura-aqui')
JWT_ALGORITHM = 'HS256'
JWT_EXPIRATION_HOURS = 24

# Credenciais fixas em memória
FIXED_USERNAME = 'admin'
FIXED_PASSWORD = 'admin123'

# URL do pipeline (para disparar a ingestão)
PIPELINE_URL = os.getenv('PIPELINE_URL', 'http://pipeline:8080/trigger') # se não existir, usa o segundo valor


def generate_token(username): # função que recebe username e retorna um token JWT para o usuário
    """Gera um token JWT para o usuário"""
    payload = { # dicionário com as informações do usuário e validade do token (payload é onde contém os dados sobre o usuário e validade do token)
        'username': username,
        'exp': datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(hours=JWT_EXPIRATION_HOURS), # momento em que o token expira
        'iat': datetime.datetime.now(datetime.timezone.utc) # momento em que foi criado o token
    }
    token = jwt.encode(payload, JWT_SECRET, algorithm=JWT_ALGORITHM) # converte o payload em um token JWT
    return token


def verify_token(f): # decorator para verificar o token JWT
    """Decorator para verificar o token JWT"""
    @wraps(f) # preserva metadados da função original, metadados são informações sobre os dados, não os dados em si
    def decorated(*args, **kwargs): # função decorada que verifica o token JWT
        token = None # define o token como None (localmente)
        
        # Verificar se o token está no header Authorization
        if 'Authorization' in request.headers: # verifica se o token está no header Authorization
            auth_header = request.headers['Authorization'] # pega o token do header Authorization
            try:
                # Formato esperado do header Authorization: "Bearer <token>", bearer indica o formato do token
                token = auth_header.split(' ')[1] # exclui "bearer" e pega apenas o token 
            except IndexError: # se o token não estiver no formato esperado
                return jsonify({'error': 'Token inválido no header'}), 401
        
        if not token: # se o token não foi fornecido
            return jsonify({'error': 'Token não fornecido'}), 401
        
        try:
            # Decodificar e verificar o token
            data = jwt.decode(token, JWT_SECRET, algorithms=[JWT_ALGORITHM]) # decodifica o token para obter os dados do usuário e validade do token
            request.current_user = data['username'] # armazena o username do usuário
        except jwt.ExpiredSignatureError: # se o token expirou
            return jsonify({'error': 'Token expirado'}), 401
        except jwt.InvalidTokenError: # se o token é inválido
            return jsonify({'error': 'Token inválido'}), 401
        
        return f(*args, **kwargs) # se o token for válido, chama a função original
    
    return decorated # retorna a função decorada


@app.route('/') # rota get para a raiz do serviço
def hello():
    return { # retorna um dicionário com as informações do serviço
        'service': 'backend1-auth',
        'status': 'running', 
        'message': 'Hello from Backend 1 (Auth & Trigger)'
    }


@app.route('/health') # rota get para o health check
def health():
    return {'status': 'healthy'}, 200


@app.route('/login', methods=['POST']) # rota post para o login, valida credenciais e gera token JWT se forem válidas
def login():
    """Endpoint de login - retorna token JWT"""
    try:
        data = request.get_json() # pega os dados do login
        
        if not data: # se os dados não forem fornecidos 
            return jsonify({'error': 'Dados não fornecidos'}), 400
        
        username = data.get('username')
        password = data.get('password')
        
        # Verificar credenciais
        if username != FIXED_USERNAME or password != FIXED_PASSWORD: # se o username ou password não forem válidos
            return jsonify({'error': 'Credenciais inválidas'}), 401
        
        # Gerar token JWT
        token = generate_token(username)
        
        return jsonify({
            'token': token,
            'username': username,
            'expires_in_hours': JWT_EXPIRATION_HOURS
        }), 200
    
    except Exception as e: # se houver erro, retorna o erro
        return jsonify({'error': f'Erro ao processar login: {str(e)}'}), 500


@app.route('/sync', methods=['POST']) # rota post para disparar o pipeline de ingestão
@verify_token # verifica o token JWT
def sync(): # função que dispara o pipeline de ingestão (sincronização)
    """Endpoint protegido para disparar o pipeline de ingestão"""
    try:
        # Fazer chamada HTTP para o pipeline
        response = requests.post( # faz uma requisição POST para o pipeline
            PIPELINE_URL,
            timeout=30  # Timeout de 30 segundos
        ) # faz uma requisição POST para o pipeline
        
        if response.status_code == 200: # se a resposta do pipeline for 200 (sucesso)
            return jsonify({
                'message': 'Pipeline de ingestão disparado com sucesso',
                'pipeline_response': response.json() if response.headers.get('content-type', '').startswith('application/json') else response.text
            }), 200
        else: # se a resposta do pipeline não for 200 (erro)
            return jsonify({
                'error': f'Erro ao disparar pipeline: status {response.status_code}', # retorna o erro
                'details': response.text # retorna o texto da resposta
            }), response.status_code # retorna o status code da resposta
    
    except requests.exceptions.Timeout: # se o timeout for excedido
        return jsonify({'error': 'Timeout ao aguardar resposta do pipeline'}), 504
    
    except requests.exceptions.ConnectionError: # se não foi possível conectar ao pipeline
        return jsonify({'error': 'Não foi possível conectar ao pipeline'}), 503
    
    except Exception as e: # se houver erro, retorna o erro
        return jsonify({'error': f'Erro ao disparar pipeline: {str(e)}'}), 500

# Executa o serviço na porta 5000 se o arquivo for executado diretamente
if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
