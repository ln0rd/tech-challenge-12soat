# Projeto Tech Challenge 12SOAT
Esse projeto se dedica as entregas da Pos Graduação em Software Architecture onde nos dedicamos a realizar esse projeto de uma oficina mecânica de médio porte, especializada em manutenção de veículos, tem enfrentado desafios para expandir seus serviços com qualidade e
eficiência.

Esse projeto tem como objetivo permitir aos clientes acompanhar em tempo real o andamento do serviço, autorizar reparos adicionais via aplicativo e garantir uma gestão interna eficiente e segura.

## Bibliotecas Utilizadas no projeto
- **HTTP Framework**: [Gorilla Mux](https://github.com/gorilla/mux) 
- **Enviroment Variables**: [godotenv](https://github.com/joho/godotenv)
- **Logger**: [zap](https://pkg.go.dev/go.uber.org/zap)

### Pré-requisitos
- **Golang Versão 1.24.5**
- **Docker** 
- **Docker Compose**

## Rodando o Projeto
Para rodar o projeto localmente utilizamos o arquivo Makefile para facilitar a execução dos comandos

Subir docker com o banco de dados
```
make up 
```
Assim que subir o container do banco de dados, podemos subir nosso projeto com:
```
make run
```
---
Podemos optar por gerar o binario e roda-lo, seguindo os comandos:
```
make build
```
E na sequencia executa-lo
```
make run-bin
```
## Estrutura Clean Architecture
Seguindo a abordagem do Clean architecture o projeto de estrutura da seguinte forma:

```
/project-root
├── /cmd                          # Entrypoints (main.go, jobs, CLI)
│   └── main.go
│
├── /internal
│   ├── /domain                   # Regras de negócio puras (Entities + Interfaces)
│   │   ├── /cliente
│   │   ├── /veiculo
│   │   ├── /insumo
│   │   └── /ordem_servico
│   │
│   ├── /usecase                  # Casos de uso (Application layer)
│   │   ├── /cliente
│   │   ├── /veiculo
│   │   ├── /insumo
│   │   └── /ordem_servico
│   │
│   ├── /interface
│   │   ├── /http                 # Handlers HTTP, rotas
│   │   └── /persistence          # Repositórios implementando interfaces de domínio
│   │
│   └── /infrastructure
│       ├── /db
│       │   └── connection.go     # Configuração do banco
│       ├── /logger
│       └── /config
│
├── /pkg                         # Bibliotecas e utilitários reutilizáveis
│   └── /utils
│
├── /migrations                  # Scripts SQL
│   ├── init_schema.sql
│   ├── add_tabela_clientes.sql
│   └── veiculos.sql
│
├── go.mod
└── README.md
```

