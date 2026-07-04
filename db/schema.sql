-- Schema do pseudo-prontuário eletrônico (espelha a especificação PSPD).
-- Serve para dois propósitos:
--   1) subir um Postgres local para desenvolvimento/teste (via docker-compose);
--   2) documentar as colunas que esperamos do banco fornecido pelo professor.
-- Se o schema real divergir, ajuste este arquivo e a camada internal/repository.

CREATE TABLE IF NOT EXISTS patients (
    patient_id  TEXT PRIMARY KEY,
    full_name   TEXT NOT NULL,
    birth_date  DATE NOT NULL,
    gender      TEXT NOT NULL,          -- male | female | other
    city        TEXT,
    state       TEXT,
    cpf         TEXT,
    cns         TEXT
);

CREATE TABLE IF NOT EXISTS encounters (
    encounter_id   TEXT PRIMARY KEY,
    patient_id     TEXT NOT NULL REFERENCES patients(patient_id),
    start_date     DATE NOT NULL,
    end_date       DATE,
    encounter_type TEXT,                -- Ambulatorial | Emergencia | Internacao | Retorno
    department     TEXT                 -- Cardiologia | Endocrinologia | Pediatria | ...
);

CREATE TABLE IF NOT EXISTS clinical_events (
    event_id     TEXT PRIMARY KEY,
    patient_id   TEXT NOT NULL REFERENCES patients(patient_id),
    encounter_id TEXT REFERENCES encounters(encounter_id),
    event_type   TEXT NOT NULL,         -- Condition | Observation | Medication
    code         TEXT NOT NULL,         -- Diabetes | Hipertensao | HbA1c | Metformina | ...
    description  TEXT,
    event_date   DATE NOT NULL,
    value        DOUBLE PRECISION,      -- preenchido quando event_type = Observation ou Medication
    unit         TEXT
);

CREATE TABLE IF NOT EXISTS user_patient_assignments (
    assignment_id       TEXT PRIMARY KEY,
    caregiver_username  TEXT NOT NULL,
    patient_id          TEXT NOT NULL REFERENCES patients(patient_id),
    assignment_type     TEXT NOT NULL,             -- medico | estagiario
    supervisor_username TEXT,
    status              TEXT NOT NULL DEFAULT 'ativo'  -- ativo | inativo
);

CREATE TABLE IF NOT EXISTS projects (
    project_id          TEXT PRIMARY KEY,
    title               TEXT NOT NULL,
    researcher_username TEXT NOT NULL,
    condition_code      TEXT NOT NULL,   -- igual ao clinical_events.code das condições
    status              TEXT NOT NULL,   -- Aprovado | Expirado | Suspenso
    valid_until         DATE
);

CREATE INDEX IF NOT EXISTS idx_encounters_patient   ON encounters(patient_id);
CREATE INDEX IF NOT EXISTS idx_events_patient       ON clinical_events(patient_id);
CREATE INDEX IF NOT EXISTS idx_events_type_code     ON clinical_events(event_type, code);
CREATE INDEX IF NOT EXISTS idx_assign_caregiver     ON user_patient_assignments(caregiver_username);
CREATE INDEX IF NOT EXISTS idx_assign_patient       ON user_patient_assignments(patient_id);
CREATE INDEX IF NOT EXISTS idx_projects_researcher  ON projects(researcher_username);
