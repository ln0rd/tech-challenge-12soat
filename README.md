# Projeto Tech Challenge 12SOAT
Esse projeto se dedica as entregas da Pos Graduação em Software Architecture onde nos dedicamos a realizar esse projeto de uma oficina mecânica de médio porte, especializada em manutenção de veículos, tem enfrentado desafios para expandir seus serviços com qualidade e
eficiência.

Esse projeto tem como objetivo permitir aos clientes acompanhar em tempo real o andamento do serviço, autorizar reparos adicionais via aplicativo e garantir uma gestão interna eficiente e segura.

## Bibliotecas Utilizadas no projeto
- **HTTP Framework**: [Gorilla Mux](https://github.com/gorilla/mux) 
- **Enviroment Variables**: [godotenv](https://github.com/joho/godotenv)
- **Logger**: [zap](https://pkg.go.dev/go.uber.org/zap)
- **GolangCi**: [golangci](github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
- **GoSec**: [gosec](https://github.com/securego/gosec)

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
### SonarQube (Análise de Código)

Para subir o SonarQube localmente com acesso via web em `http://localhost:9000`:

1. Subir os serviços do SonarQube:
   ```
   make sonar-up
   ```
2. Acessar a interface web: `http://localhost:9000` (credenciais padrão: `admin` / `admin`). Ao logar pela primeira vez, será solicitado alterar a senha.
3. Criar um token de usuário no seu perfil do SonarQube e exportá-lo localmente:
   ```
   export SONAR_TOKEN=your_token
   ```
4. Executar a análise local (requer `sonar-scanner` instalado):
   ```
   make sonar-scan
   ```
   
   Ou usar o comando simples sem token (para desenvolvimento local):
   ```
   make sonar-scan
   ```
   
   O comando `sonar-scan` automaticamente gera a cobertura de testes e relatórios de segurança antes de executar a análise.
   
   **Nota:** Para Security e Reliability, o SonarQube usa o `golangci-lint` com regras de segurança habilitadas (incluindo `gosec`).

Senha padrão do sonarQube admin/admin;
Sugestão de senha para testes no sonar: Senhaforte123@

Arquivo de configuração `sonar-project.properties` já incluído na raiz do projeto.
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

