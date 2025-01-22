# Projeto Cotação de Dólar

Este projeto contém um **servidor** que expõe uma API de cotação de dólar e um **cliente** em Go que consome essa API.

## Requisitos

- **Docker** para rodar o servidor via Docker Compose.
- **Go** para rodar o cliente.

---

## Passo 1: Rodando o Servidor

O servidor é configurado para rodar via **Docker Compose** na porta **8080** e oferece um endpoint para obter a cotação do dólar. O endpoint é `/cotacao`.

1. Clone o repositório:
   ` git clone https://github.com/HenriqueOtsuka/dolar-price.git `

2. Navegue até a pasta do projeto:
   ` cd server `

3. Execute o servidor usando **Docker Compose**:
   ` docker-compose up --build `

Isso irá iniciar o servidor na porta **8080**, que estará disponível em `http://localhost:8080/cotacao`.

---

## Passo 2: Rodando o Cliente

O cliente é escrito em **Go** e irá fazer uma requisição HTTP para o servidor e exibir a cotação do dólar.

Para rodar o cliente, execute o seguinte comando:
   ` go run client.go `

Isso fará com que o cliente consuma o endpoint `/cotacao` do servidor e exiba o preço do dólar.

---

## Passo 3: Como Funciona

### Servidor:

O servidor expõe um endpoint `GET /cotacao`, que ao ser acessado, retorna a cotação do dólar (USD para BRL) em formato JSON. A resposta será similar a:

   ` 
   {
       "bid": "6.30",
       "ask": "6.35"
   }
   `

(Mas eu estou trazendo tudo)

Além disso, ele insere no sqlite que está dentro do container Docker. `docker exec -it container_id /bin/sh`

Depois disso é só dar um `ls` que vai ser possível ver o database.db

```
.tables;

.headers ON;

SELECT * FROM dolar_price;
```

### Cliente:

O cliente faz uma requisição HTTP para o servidor e exibe o valor da cotação, retornando algo como:

   ` 
   Preço do Dólar: 6.30
   `

Quando ele roda, escreve a resposta da cotação atual em um arquivo chamado cotacao.txt que vai ser criado no mesmo diretório do client.

---
