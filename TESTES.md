# Guia de Testes — patient-data-service

Como subir o serviço e validar todas as consultas e agregações **pelo Postman** (gRPC) **contra o
banco oficial do professor** (`pseudopep_g03`), com os payloads prontos para copiar e colar e os
resultados reais esperados.

> No Postman (gRPC) os campos usam os **nomes do `.proto` (snake_case)**, ex.: `patient_id`.
> Os códigos/valores seguem a convenção do banco real (MAIÚSCULO/inglês): `DIABETES`, `HBA1C`,
> `OBSERVATION`, `ATTENDING`, etc.

---

## 1. Pré-requisitos
- **Go 1.25+**
- **Postman** (versão com suporte a gRPC)
- **Acesso ao banco do professor** via túnel SSH (ou, offline, Docker + Compose para o Postgres local)

## 2. Subir o serviço (contra o banco do professor)
O banco `pseudopep_g03` fica na rede interna do cluster; a partir da sua máquina, alcança-se por um
túnel SSH.

```bash
# janela 1 — túnel SSH (deixe aberta)
ssh -p 10200 -L 15432:192.168.122.1:5432 <matricula>@kiriland.unb.br

# janela 2 — serviço apontando para o túnel
cp .env.example .env     # DATABASE_URL no formato keyword/value do pgx (a senha tem '@'):
                         # host=localhost port=15432 user=grupo03_user password=123@g03
                         # dbname=pseudopep_g03 sslmode=disable
go run ./cmd/server
```
Logs esperados: `conectado ao Postgres`, `servidor gRPC ouvindo port=50051`.
- gRPC: **localhost:50051** · Métricas/health: **http://localhost:9090**

> **Confirme que é o banco certo:** rode o CT-11 (`GetCohortStatistics DIABETES`) — deve voltar
> `total 30.110`. O seed local daria 13.
>
> **Offline (seed local):** `docker compose up -d db` + `.env` apontando para `localhost:5433`
> (ou tudo em container: `docker compose up -d --build`).

## 3. Testes automatizados
```bash
go test ./...                    # unitários (sem banco)
```
> Os testes de **integração** (`go test -tags=integration ./...`) checam as **invariantes do seed
> local** (coorte de 13, `med.cardoso` com 8, etc.) — rode-os contra o **Postgres local**
> (`docker compose up -d db`), **não** contra o banco do professor (lá a coorte tem 30.110 e os
> números não batem com o seed).

---

## 4. Configurar o Postman (uma vez)
1. **New → gRPC Request**.
2. Em **Enter server URL**: `localhost:50051` (deixe **TLS desligado**).
3. O serviço tem **server reflection**: o Postman lista os métodos sozinho no dropdown **Select a method**.
   (Alternativa: **Import a .proto file** → `proto/patientdata/v1/patientdata.proto`.)
4. Escolha o método, cole o JSON na aba **Message** e clique **Invoke**.

> Todos os métodos ficam sob `patientdata.v1.PatientDataService`.

---

## 5. Casos de teste (contra o banco do professor)

### CT-01 — GetPatient (dados completos)
Esperado: **Ana Almeida**, com CPF/CNS, nascimento 2009-03-12, Formosa/GO.
```json
{ "patient_id": "P030000001" }
```

### CT-02 — GetPatient (inexistente → erro)
Esperado: erro `NOT_FOUND`.
```json
{ "patient_id": "P999999" }
```

### CT-03 — ListPatientsByDoctor  *(server streaming)*
Esperado: **stream** com os pacientes vinculados ao médico `med.almeida` (~**30.001** — a base é grande); o grpcurl/Postman mostram os `Patient` chegando um a um.
```json
{ "doctor_username": "med.almeida" }
```

### CT-04 — ListSupervisedPatients  *(server streaming)*
Esperado: **stream** com os pacientes supervisionados pelo estagiário `est.ferreira` (~**7.447**).
```json
{ "intern_username": "est.ferreira" }
```

### CT-05 — ListEncounters
Esperado: os atendimentos do paciente — ex.: `E03000000001` (INPATIENT, PEDIATRICS, 2024-08-05).
```json
{ "patient_id": "P030000001" }
```

### CT-06 — ListClinicalEvents (só exames)
Esperado: os eventos do tipo `OBSERVATION` do paciente com valor e unidade — ex.: `BMI` 26,9 kg/m².
```json
{ "patient_id": "P030000001", "event_type": "OBSERVATION" }
```

### CT-07 — ListClinicalEvents (todos os eventos)
Esperado: condições + exames + medicações do paciente.
```json
{ "patient_id": "P030000001" }
```

### CT-08 — GetClinicalSummary
Esperado: Ana Almeida + último atendimento (PEDIATRICS) + condição `HEART_FAILURE` + exame `BMI` 26,9 + medicação `SIMVASTATIN` 20 mg.
```json
{ "patient_id": "P030000001" }
```

### CT-09 — GetClinicalHistory
Esperado: eventos em ordem temporal crescente.
```json
{ "patient_id": "P030000001" }
```

### CT-10 — ListCohortPatients  *(server streaming)*
Esperado: **stream** com **30.110** pacientes diabéticos, chegando um a um (sem bufferizar).
```json
{ "condition_code": "DIABETES" }
```

### CT-11 — GetCohortStatistics (coorte real de Diabetes)
Esperado (dados reais):
- total **30.110**;
- sexo F **15.022** (49,9%) / M **15.088** (50,1%);
- faixas 0-17 **4.613** · 18-39 **8.145** · 40-59 **7.362** · 60+ **9.990** (somam o total);
- HbA1c média **8,50** / mediana **8,5**;
- medicamentos LOSARTAN **11.362** · ENALAPRIL **11.354** · SIMVASTATIN **11.206** · METFORMIN **11.139** · INSULIN **11.084**;
- **14 departamentos** (PEDIATRICS 5.107 … GERIATRICS 4.876, ~5 mil cada).
```json
{ "condition_code": "DIABETES" }
```

### CT-12 — ListProjectsByResearcher
Esperado: **2** projetos do `pes.mendes` — PRJ01_G03 (DIABETES, APPROVED) e PRJ04_G03 (CKD, PENDING).
```json
{ "researcher_username": "pes.mendes" }
```

### CT-13 — CheckAssignment (vínculo válido)
Esperado: `allowed: true`, `assignment_type: "ATTENDING"` (o `role` aceita `medico`/`estagiario` e é
mapeado). O par `med.almeida` ↔ `P030000001` existe no banco.
```json
{ "username": "med.almeida", "patient_id": "P030000001", "role": "medico" }
```

### CT-14 — CheckAssignment (sem vínculo)
Esperado: `allowed: false` (paciente inexistente → sem vínculo).
```json
{ "username": "med.almeida", "patient_id": "P999999", "role": "medico" }
```

---

## 6. Observabilidade (Prometheus)
Depois de rodar alguns casos, abra no navegador ou via curl:
- `http://localhost:9090/healthz` → `ok`
- `http://localhost:9090/metrics` → procure por `db_queries_total`, `grpc_server_handled_total`,
  `grpc_server_handling_seconds`, `pgxpool_total_conns`, `go_goroutines`, `process_resident_memory_bytes`.

## 7. Encerrar
Pare o serviço com `Ctrl+C` e feche a janela do túnel SSH. (No modo offline, `docker compose down -v`.)

---

## Referência — dados reais do banco do professor (`pseudopep_g03`)

**Usuários / vínculos (`user_patient_assignments`):**
| Papel | Usuários (exemplos) |
|---|---|
| Médico (`ATTENDING`) | `med.rocha` (30.092), `med.monteiro` (30.036), `med.lima` (30.025), `med.almeida` (30.001), `med.cardoso` (29.846) |
| Estagiário (`TRAINEE`) | `est.dias` (7.629), `est.melo` (7.596), `est.costa` (7.496), `est.ferreira` (7.447), `est.gomes` (7.410) |
| Pesquisador (`projects`) | `pes.mendes` (PRJ01_G03 DIABETES/APPROVED, PRJ04_G03 CKD/PENDING); `pes.araujo` (PRJ02_G03 HYPERTENSION/APPROVED, PRJ05_G03 HEART_FAILURE/EXPIRED); `pes.silveira` (PRJ03_G03 OBESITY/APPROVED) |

Exemplo de vínculo: `med.almeida` ↔ `P030000001` (ATTENDING) · `est.ferreira` ↔ `P030000001` (TRAINEE).

**Coorte DIABETES:**
| Métrica | Valor real |
|---|---|
| Total | 30.110 |
| Sexo | F 15.022 (49,9%) / M 15.088 (50,1%) |
| Faixa etária | 0-17: 4.613 · 18-39: 8.145 · 40-59: 7.362 · 60+: 9.990 |
| HbA1c | média 8,50 / mediana 8,5 |
| Medicamentos | LOSARTAN 11.362 · ENALAPRIL 11.354 · SIMVASTATIN 11.206 · METFORMIN 11.139 · INSULIN 11.084 |
| Departamentos | 14 (PEDIATRICS 5.107 … GERIATRICS 4.876) |

Paciente real de exemplo: **P030000001** (Ana Almeida, 2009-03-12, female, Formosa/GO).

> Convenções do banco real: MAIÚSCULO/inglês, `assignment_type` `ATTENDING`/`TRAINEE`, `active`
> booleano, `target_condition_code`, `value` como texto. O seed local (`db/seed.sql`) espelha essas
> convenções em escala menor (coorte de 13) para os testes de integração automatizados.

## Troubleshooting
| Sintoma | Solução |
|---|---|
| `conectado ao Postgres` não aparece / timeout ao subir | o túnel SSH não está aberto ou a porta está errada; confira a janela 1 (`-L 15432:192.168.122.1:5432`) |
| `password authentication failed` | a senha tem `@`; use o formato **keyword/value** do pgx no `DATABASE_URL`, não a URL `postgres://` |
| `GetCohortStatistics DIABETES` volta **13** em vez de 30.110 | você está no **seed local**, não no banco do professor; confira o `DATABASE_URL`/túnel |
| `patient_id é obrigatório` (InvalidArgument) | o campo chegou vazio; confira o **nome snake_case** (`patient_id`, não `patientId`) e o valor |
| Postman não conecta | serviço fora do ar ou porta errada (gRPC é 50051); confira que o TLS está desligado |
| métodos não aparecem no Postman | a reflection precisa do serviço no ar; recarregue após conectar, ou importe o `.proto` |
