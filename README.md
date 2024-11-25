# Concorrência com Golang - Leilão
Este projeto demonstra o fechamento automático de leilões usando concorrência em Go, baseado em um tempo definido pelo usuário.

## Configuração e Execução
**1. Configurar o Tempo de Fechamento**
* No arquivo `.env`, ajuste o valor do campo `AUCTION_DURATION` para definir o tempo que um leilão permanecerá ativo antes de ser fechado automaticamente.

**2. Iniciar os Serviços**
* Execute o comando abaixo para iniciar a aplicação:
    ```
    docker-compose up -d
    ```
* Aguarde até que todos os serviços estejam completamente inicializados.
## Testando o Fechamento Automático
**1. Criar um Novo Leilão**

 Faça uma requisição POST para criar um leilão. Use o exemplo abaixo como modelo:
```http
POST http://localhost:8080/auction HTTP/1.1
Content-Type: application/json

{
    "product_name": "celular",
    "category": "eletronicos",
    "description": "Iphone 1 8GB",
    "condition": 0
}
```
**2. Listar Leilões Ativos**

 Use a seguinte requisição GET para listar leilões com o status ativo (status = 0):
```http
GET http://localhost:8080/auction?status=0 HTTP/1.1
```
* O campo status dos leilões retornados será 0, indicando que estão ativos.

**3. Verificar o Fechamento**

* Aguarde o tempo configurado em `AUCTION_DURATION` no `.env`.
* Após o período, repita a requisição GET para listar os leilões:
    ```http
    GET http://localhost:8080/auction?status=1 HTTP/1.1
    ```
* Agora, o campo status dos leilões será 1, indicando que foram fechados automaticamente.