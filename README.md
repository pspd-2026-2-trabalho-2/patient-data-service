# patient-data-service

Microsserviço **gRPC em Go** do projeto PSPD (Trabalho 2 — observabilidade de
microsserviços em Kubernetes). É a **camada de dados** do backend: executa
consultas SQL sobre o pseudo-prontuário eletrônico do Hospital Universitário e
devolve **dados crus** (sem anonimização).

A aplicação das regras de acesso (FULL/PARTIAL/ANONYMIZED/AGGREGATED) e a
conversão para HL7/FHIR são responsabilidade do **data-transform-service**
(repositório separado), que consome este serviço.

```
API Gateway ──(gRPC)──▶ patient-data-service ──(SQL)──▶ PostgreSQL
                                   │
                                   └── expõe /metrics (Prometheus) e health check
```

## Responsabilidades

- Consultas: pacientes por médico, supervisionados por estagiário, atendimentos,
  eventos clínicos (condições/exames/medicações), resumo e histórico clínico,
  pacientes de uma coorte e projetos de pesquisa.
- Agregações da coorte: total, distribuição por sexo, faixa etária e departamento,
  média e mediana de HbA1c, frequência de medicamentos.
- `CheckAssignment`: verifica o vínculo cuidador↔paciente (usado pelo
  Authorization Service — este serviço é o único dono do banco).
- Observabilidade: métricas Prometheus (RPC, consultas ao banco, pool, Go/processo).

## Requisitos

- Go 1.25+
- Docker + Docker Compose (para o Postgres local)
- `protoc` + plugins Go (só para **regerar** o código do `.proto`; o código já vem
  gerado em `gen/`)
- Opcional: `grpcurl` para testar manualmente

## Estrutura

```
proto/patientdata/v1/patientdata.proto   contrato gRPC (fonte de verdade)
gen/patientdata/v1/                       código gerado (não editar)
db/schema.sql, db/seed.sql                schema + dados de exemplo
cmd/server/main.go                        entrypoint (wiring + shutdown)
internal/config                           variáveis de ambiente
internal/db                               pool pgx
internal/domain                           modelos de negócio
internal/repository                       SQL (uma função por consulta)
internal/service                          regras + agregações
internal/grpcserver                       transporte gRPC (domain↔protobuf)
internal/observability                    métricas Prometheus
```

## Como executar

### Opção A — contra o banco do professor (via túnel SSH) — recomendado

A base oficial `pseudopep_g03` fica na rede interna do cluster; a partir da sua máquina, alcança-se
por um túnel SSH:

```bash
# janela 1 — túnel (deixe aberta)
ssh -p 10200 -L 15432:192.168.122.1:5432 <matricula>@kiriland.unb.br

# janela 2 — serviço apontando para o túnel (ver .env)
cp .env.example .env             # DATABASE_URL = host=localhost port=15432 ... dbname=pseudopep_g03
go run ./cmd/server
```

### Opção B — Postgres local (desenvolvimento offline)

```bash
docker compose up -d db          # sobe o Postgres local (schema + seed)
cp .env.example .env             # aponte DATABASE_URL para localhost:5433
go run ./cmd/server
# ou tudo em container: docker compose up -d --build
```

O serviço sobe em:
- gRPC: `localhost:50051`
- Métricas/health HTTP: `localhost:9090` (`/metrics`, `/healthz`)

## Variáveis de ambiente

| Variável       | Padrão                                                         | Descrição                       |
|----------------|---------------------------------------------------------------|---------------------------------|
| `GRPC_PORT`    | `50051`                                                       | porta do servidor gRPC          |
| `METRICS_PORT` | `9090`                                                        | porta HTTP de métricas/health   |
| `DATABASE_URL` | `postgres://pspd:pspd@localhost:5433/hospital?sslmode=disable`| conexão (prioritária)           |
| `PG*`          | ver `.env.example`                                            | usadas se `DATABASE_URL` vazia  |
| `LOG_LEVEL`    | `info`                                                        | `debug`/`info`/`warn`/`error`   |

O serviço já foi validado contra o banco oficial do professor (`pseudopep_g03`): basta apontar
`DATABASE_URL` para ele — sem mudar código. Como a senha tem `@`, use o **formato keyword/value**
do pgx (não a URL `postgres://`):

```
DATABASE_URL="host=localhost port=15432 user=grupo03_user password=123@g03 dbname=pseudopep_g03 sslmode=disable"
```

## RPCs

| RPC | Entrada | Saída |
|-----|---------|-------|
| `ListPatientsByDoctor` | `doctor_username` | lista de pacientes |
| `ListSupervisedPatients` | `intern_username` | lista de pacientes |
| `GetPatient` | `patient_id` | paciente |
| `ListEncounters` | `patient_id` | atendimentos |
| `ListClinicalEvents` | `patient_id`, `event_type?` | eventos clínicos |
| `GetClinicalSummary` | `patient_id` | resumo clínico |
| `GetClinicalHistory` | `patient_id` | histórico (ordem temporal) |
| `ListCohortPatients` | `condition_code` | pacientes da coorte |
| `GetCohortStatistics` | `condition_code` | agregações |
| `ListProjectsByResearcher` | `researcher_username` | projetos |
| `CheckAssignment` | `username`, `patient_id`, `role?` | `allowed`, `assignment_type` |

## Testes

Guia completo (como subir, casos de teste, payloads e resultados esperados) em
**[TESTES.md](TESTES.md)**. Resumo rápido:

```bash
go test ./...                       # unitários (não precisam de banco)
docker compose up -d db             # sobe o Postgres local com seed
go test -tags=integration ./...     # integração (valida invariantes do SEED local)
```
> A integração checa as invariantes do **seed local** (coorte de 13). A validação com **dados reais**
> (coorte de 30.110 no banco do professor) é feita manualmente via Postman/grpcurl — ver **TESTES.md**.

## Observabilidade

`GET http://localhost:9090/metrics` expõe, entre outras:

- `grpc_server_handled_total{method,code}` e `grpc_server_handling_seconds` — RPCs e latência
- `db_queries_total{query,status}`, `db_query_duration_seconds`, `db_rows_returned` — consultas ao banco
- `pgxpool_*` — conexões do pool
- `go_*` / `process_*` — CPU, memória, goroutines

## Regenerar o código do proto

```bash
make tools    # instala protoc-gen-go e protoc-gen-go-grpc
make proto    # ou o comando protoc equivalente (ver Makefile)
```
