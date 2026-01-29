# Dockerfile na raiz para deploy no Fly.io
# Este Dockerfile é um wrapper para o frontend
FROM node:18-alpine

WORKDIR /app

# Copiar arquivos de dependências do frontend
COPY frontend/package.json frontend/package-lock.json* ./

# Instalar dependências
RUN npm install

# Copiar código fonte do frontend
COPY frontend/ .

# Build da aplicação React
RUN npm run build

# Instalar serve para servir os arquivos estáticos
RUN npm install -g serve

# Expor porta
EXPOSE 3000

# Comando para servir os arquivos estáticos
CMD ["serve", "-s", "build", "-l", "3000"]

