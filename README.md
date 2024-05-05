# Execução e testes

## Ambiente DEV com uso de docker compose
Para subir a imagem fazendo o build do projeto, executar o comendo abaixo" :

    docker compose up --build

A execução do comando acima irá subir uma imagem com os seguintes serviços:
    
    - jaeger
    - zipkin
    - prometheus
    - otel-collector
    
e irá dispnibilizar o código do app deste projeto na pasta: 

    src/app/

Para entrar na imagem para desenvolvimento, executar o comando abaixo:

    docker compose exec goapp bash

Para rodar o sistema em dev, estando na pasta citada acima, executar o comando abaixo:

    go run main.go 

Após a execução do comando acima, os serviços do app estarão disponíveis na porta 8080. <br>
Para realizar os testes, basta fazer uma chamada POST passando o cep a ser coletado da seguinte forma:

    POST http://localhost:8080/ HTTP/1.1
    content-type: application/json

    {
        "cep": "04112080"
    }

O projeto possui o arquivo request.http que pode ser utilizado para os testes.<BR>
A chamda acima, irá iniciar o processamanto pelo "serviço A", que írá relaizar as validações necessárias do CEP passado como parâmetro. Caso haja algum problema de validação, isso será informado conforme requisitos do projeto.
<BR>
Após as validações necessárias do CEP, o "serviço A" por sua vez irá chamar o "serviço B" na seguinte URL:

    http://localhost:8080/clima?cep={CEP recebido como parâmetro no serviço A}

O serviço B írá relaizar as validações necessárias do CEP passado como parâmetro. Caso haja algum problema de validação, isso será informado conforme requisitos do projeto.

Caso estaja tudo certo com o CEP, o "serviço B" irá coletar as informações de localidade do CEP. Caso a consulta tenha sucesso, o "serviço B" irá coletar as informações de temperatura atual da localidade. Caso obtenha sucesso irá retornar as informações de temperatura no seguinte formato:

    {
    "city": "São Paulo",
    "temp_C": 29,
    "temp_F": 84.2,
    "temp_K": 302
    }

Para ambas as chamdas de serviços externos que o "serviço B" realiza, caso ocorra algum erro, isto será informado conforme requisitos do projeto.


# Observabilidade 

## Jaeger
Acessar o serviço no endereço:

    http://localhost:16686/search

Após a execução de uma consulta ao "serviço A" <B>com sucesso</B>, na combo "Sevice" deverá estar diponível o item

    otel_Lab

Selecionando este item e clicando em "Find Traces" deverá ser encontrado ao menos 1 "Trace" com 3 spans. Abrindo esse Trace poderá ser vista toda a rastrabilidade da execução dos spans.

## zipkin

Acessar o serviço no endereço:

    http://localhost:9411/zipkin/

Após a execução de uma consulta ao "serviço A" <B>com sucesso</B>, clicar em "Run query" na tela principal do zipkin.

Serão exibidas as rastreabilidades onde a principal delas é 

    consulta temperaturas: otel_lab

Entrando nos detalhes deste item será exibida a rastreabilidade completa de cada serviço executado

