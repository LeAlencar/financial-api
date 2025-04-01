# filepath: e:\Projetos\financial-api\Dockerfile
# Usar uma imagem base do Go
FROM golang:1.20-alpine

# Configurar o diretório de trabalho
WORKDIR /financial-api
# Definir variáveis de ambiente


# Copiar os arquivos do projeto para o contêiner
COPY . .

# Baixar as dependências
RUN go mod tidy

# Compilar o binário
RUN go build -o main .

# Expor a porta 8080
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]