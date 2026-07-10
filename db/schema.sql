-- Schema do pseudoprontuário eletrônico, alinhado ao banco REAL fornecido pelo
-- professor (pseudopep_gNN). Serve para subir um Postgres local de desenvolvimento
-- com as MESMAS convenções do banco oficial (colunas, tipos e valores em MAIÚSCULO/inglês).

CREATE TABLE IF NOT EXISTS patients (
    patient_id  VARCHAR(20) PRIMARY KEY,
    full_name   VARCHAR(120) NOT NULL,
    birth_date  DATE NOT NULL,
    gender      VARCHAR(20) NOT NULL CHECK (gender IN ('male','female')),
    city        VARCHAR(80) NOT NULL,
    state       CHAR(2) NOT NULL,
    cpf         VARCHAR(14) NOT NULL,
    cns         VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS encounters (
    encounter_id   VARCHAR(20) PRIMARY KEY,
    patient_id     VARCHAR(20) NOT NULL REFERENCES patients(patient_id),
    start_date     TIMESTAMP NOT NULL,
    end_date       TIMESTAMP,
    encounter_type VARCHAR(40) NOT NULL,  -- AMBULATORIAL | EMERGENCY | INPATIENT | ICU | FOLLOW_UP | TELEHEALTH
    department     VARCHAR(80) NOT NULL   -- CARDIOLOGY | ENDOCRINOLOGY | NEPHROLOGY | ...
);

CREATE TABLE IF NOT EXISTS clinical_events (
    event_id     VARCHAR(20) PRIMARY KEY,
    patient_id   VARCHAR(20) NOT NULL REFERENCES patients(patient_id),
    encounter_id VARCHAR(20) NOT NULL REFERENCES encounters(encounter_id),
    event_type   VARCHAR(30) NOT NULL CHECK (event_type IN ('CONDITION','OBSERVATION','MEDICATION')),
    code         VARCHAR(40) NOT NULL,   -- DIABETES | HYPERTENSION | HBA1C | GLUCOSE | METFORMIN | ...
    description  VARCHAR(200) NOT NULL,
    value        VARCHAR(80),            -- texto; numérico para OBSERVATION/MEDICATION
    unit         VARCHAR(40),
    event_date   TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS user_patient_assignments (
    assignment_id       VARCHAR(20) PRIMARY KEY,
    username            VARCHAR(80) NOT NULL,
    patient_id          VARCHAR(20) NOT NULL REFERENCES patients(patient_id),
    assignment_type     VARCHAR(20) NOT NULL CHECK (assignment_type IN ('ATTENDING','TRAINEE')),
    supervisor_username VARCHAR(80),
    active              BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS projects (
    project_id            VARCHAR(20) PRIMARY KEY,
    title                 VARCHAR(200) NOT NULL,
    researcher_username   VARCHAR(80) NOT NULL,
    target_condition_code VARCHAR(40) NOT NULL,   -- igual ao clinical_events.code das condições
    status                VARCHAR(40) NOT NULL,    -- APPROVED | PENDING | EXPIRED | REJECTED | SUSPENDED
    valid_until           DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_encounters_patient   ON encounters(patient_id);
CREATE INDEX IF NOT EXISTS idx_events_patient        ON clinical_events(patient_id);
CREATE INDEX IF NOT EXISTS idx_events_type_code      ON clinical_events(event_type, code);
CREATE INDEX IF NOT EXISTS idx_assign_username       ON user_patient_assignments(username);
CREATE INDEX IF NOT EXISTS idx_assign_patient        ON user_patient_assignments(patient_id);
CREATE INDEX IF NOT EXISTS idx_projects_researcher   ON projects(researcher_username);
