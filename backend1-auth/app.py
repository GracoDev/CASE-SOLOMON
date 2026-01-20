from flask import Flask

app = Flask(__name__)

@app.route('/')
def hello():
    return {'service': 'backend1-auth', 'status': 'running', 'message': 'Hello from Backend 1 (Auth & Trigger)'}

@app.route('/health')
def health():
    return {'status': 'healthy'}, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)



