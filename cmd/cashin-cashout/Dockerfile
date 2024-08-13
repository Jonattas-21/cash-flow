# Usando a imagem oficial do Go como base
FROM golang:1.22-alpine

# Criando o diretório de trabalho
WORKDIR /app

# Copiando os arquivos do serviço para o diretório de trabalho
COPY . .

# Baixando as dependências e construindo o binário
RUN go mod tidy
RUN go build -o cashin-cashout ./cmd/cashin-cashout/main.go

# Comando de inicialização do container
CMD ["./cashin-cashout"]