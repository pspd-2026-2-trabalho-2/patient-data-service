# Guia de Testes — patient-data-service

Como subir o serviço e validar todas as consultas e agregações **pelo Postman** (gRPC), com os
payloads prontos para copiar e colar e os resultados esperados (baseados no `db/seed.sql`).

> No Postman (gRPC) os campos usam os **nomes do `.proto` (snake_case)**, ex.: `patient_id`.

---

## 1. Pré-requisitos
- **Go 1.25+** e **Docker + Docker Compose**
- **Postman** (versão com suporte a gRPC — qualquer versão recente)

## 2. Subir o serviço
Na raiz do repositório:
```bash
docker compose up -d db          # Postgres já com schema + seed (porta 5433 no host)
cp .env.example .env             # aponta DATABASE_URL para o Postgres local
go run ./cmd/server              # sobe o serviço
```
Logs esperados: `conectado ao Postgres`, `servidor gRPC ouvindo port=50051`.
- gRPC: **localhost:50051** · Métricas/health: **http://localhost:9090**

> Alternativa (tudo em container): `docker compose up -d --build`.

## 3. Testes automatizados (opcional, mas recomendado)
```bash
go test ./...                    # unitários (sem banco)
go test -tags=integration ./...  # integração: confere os números do seed automaticamente
```

---

## 4. Configurar o Postman (uma vez)
1. **New → gRPC Request**.
2. Em **Enter server URL**: `localhost:50051` (deixe **TLS desligado**).
3. O serviço tem **server reflection**: o Postman lista os métodos sozinho no dropdown **Select a method**.
   (Alternativa: **Import a .proto file** → `proto/patientdata/v1/patientdata.proto`.)
4. Escolha o método, cole o JSON na aba **Message** e clique **Invoke**.

> Todos os métodos ficam sob `patientdata.v1.PatientDataService`.

---

## 5. Casos de teste

### CT-01 — GetPatient (dados completos)
Esperado: João da Silva, com CPF/CNS, nascimento 1970-05-10, Brasilia/DF.
```json
{ "patient_id": "P000001" }
```

### CT-02 — GetPatient (inexistente → erro)
Esperado: erro `NOT_FOUND`.
```json
{ "patient_id": "P999999" }
```

### CT-03 — ListPatientsByDoctor
Esperado: **8** pacientes (P000001–P000008).
```json
{ "doctor_username": "med.cardoso" }
```

### CT-04 — ListSupervisedPatients
Esperado: **3** pacientes (P000001–P000003).
```json
{ "intern_username": "est.souza" }
```

### CT-05 — ListEncounters
Esperado: **2** atendimentos (ENC16 Cardiologia 2024-08-01, ENC01 Endocrinologia 2024-02-10).
```json
{ "patient_id": "P000001" }
```

### CT-06 — ListClinicalEvents (só exames)
Esperado: HBA1C 8.1 % e GLUCOSE 182 mg/dL.
```json
{ "patient_id": "P000001", "event_type": "OBSERVATION" }
```

### CT-07 — ListClinicalEvents (todos os eventos)
Esperado: condições + exames + medicações do paciente.
```json
{ "patient_id": "P000001" }
```

### CT-08 — GetClinicalSummary
Esperado: paciente + último atendimento + condições + exames + medicações.
```json
{ "patient_id": "P000001" }
```

### CT-09 — GetClinicalHistory
Esperado: eventos em ordem temporal crescente.
```json
{ "patient_id": "P000001" }
```

### CT-10 — ListCohortPatients
Esperado: **13** pacientes diabéticos.
```json
{ "condition_code": "DIABETES" }
```

### CT-11 — GetCohortStatistics
Esperado: total 13; sexo F7/M6; faixas 3/6/4; HbA1c média 7.76 / mediana 7.6;
medicamentos METFORMIN 10, INSULIN 4, LOSARTAN 2; departamentos ENDOCRINOLOGY 13, CARDIOLOGY 2.
```json
{ "condition_code": "DIABETES" }
```

### CT-12 — ListProjectsByResearcher
Esperado: **2** projetos (PRJ01 APPROVED, PRJ02 EXPIRED).
```json
{ "researcher_username": "pesq.lima" }
```

### CT-13 — CheckAssignment (vínculo válido)
Esperado: `allowed: true`, `assignment_type: "ATTENDING"` (o `role` aceita `medico`/`estagiario` e é mapeado).
```json
{ "username": "med.cardoso", "patient_id": "P000001", "role": "medico" }
```

### CT-14 — CheckAssignment (sem vínculo)
Esperado: `allowed: false` (P000015 é do med.almeida).
```json
{ "username": "med.cardoso", "patient_id": "P000015", "role": "medico" }
```

---

## 6. Observabilidade (Prometheus)
Depois de rodar alguns casos, abra no navegador ou via curl:
- `http://localhost:9090/healthz` → `ok`
- `http://localhost:9090/metrics` → procure por `db_queries_total`, `grpc_server_handled_total`,
  `grpc_server_handling_seconds`, `pgxpool_total_conns`, `go_goroutines`, `process_resident_memory_bytes`.

## 7. Encerrar
Pare o serviço com `Ctrl+C` e depois `docker compose down -v` (remove o Postgres e o volume).

---

## Referência — dados do seed
| Usuário | Papel | Pacientes |
|---|---|---|
| `med.cardoso` | médico | P000001–P000008 (8) |
| `est.souza` | estagiário (sup.: med.cardoso) | P000001–P000003 (3) |
| `med.almeida` | médico | P000009–P000015 (7) |
| `pesq.lima` | pesquisador | PRJ01 (DIABETES/APPROVED), PRJ02 (HYPERTENSION/EXPIRED) |

**Coorte DIABETES (13):** P000001–P000012 e P000015 · sexo 7F/6M · faixas 18-39=3, 40-59=6, 60+=4 ·
HbA1c média 7.76 / mediana 7.6 · METFORMIN 10, INSULIN 4, LOSARTAN 2 · ENDOCRINOLOGY 13, CARDIOLOGY 2.

> Estas convenções (MAIÚSCULO/inglês, `assignment_type` ATTENDING/TRAINEE, `active` booleano,
> `target_condition_code`, `value` como texto) são as mesmas do banco real do professor
> (`pseudopep_gNN`). Para testar contra ele, veja o `.env` (conexão via túnel SSH).

## Troubleshooting
| Sintoma | Solução |
|---|---|
| `patient_id é obrigatório` (InvalidArgument) | o campo chegou vazio; confira o **nome snake_case** (`patient_id`, não `patientId`) e que o valor está preenchido |
| Postman não conecta | serviço fora do ar ou porta errada (gRPC é 50051); confira que o TLS está desligado |
| `ping no banco` falha ao subir | Postgres não pronto; `docker compose up -d db` e aguarde `healthy` |
| porta 5433 em uso | outro Postgres no host; ajuste a porta no `docker-compose.yml` e no `.env` |
| métodos não aparecem no Postman | a reflection precisa do serviço no ar; recarregue após conectar, ou importe o `.proto` |
