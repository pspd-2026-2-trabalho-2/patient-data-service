-- Dados de exemplo para desenvolvimento/teste sem depender do banco do professor.
-- Coorte "Diabetes" projetada para dar agregações verificáveis:
--   13 pacientes diabéticos (P000001..P000012 e P000015)
--   por sexo: 6 male / 7 female
--   por faixa etária (ref. 2026): 18-39 = 3, 40-59 = 6, 60+ = 4
--   HbA1c: média ~7.76, mediana 7.6
--   medicamentos na coorte: Metformina=10, Insulina=4, Losartana=2
-- P000013 e P000014 NÃO são diabéticos (só Hipertensão) -> validam o filtro de coorte.

-- ---------- patients ----------
INSERT INTO patients (patient_id, full_name, birth_date, gender, city, state, cpf, cns) VALUES
  ('P000001','João da Silva',      '1970-05-10','male',  'Brasilia','DF','111.111.111-11','700000000000001'),
  ('P000002','Maria Souza',        '1985-03-22','female','Goiania', 'GO','222.222.222-22','700000000000002'),
  ('P000003','Ana Oliveira',       '1960-11-02','female','Brasilia','DF','333.333.333-33','700000000000003'),
  ('P000004','Carlos Pereira',     '1978-07-15','male',  'Anapolis','GO','444.444.444-44','700000000000004'),
  ('P000005','Beatriz Santos',     '1995-01-30','female','Brasilia','DF','555.555.555-55','700000000000005'),
  ('P000006','Pedro Costa',        '1955-09-09','male',  'Luziania','GO','666.666.666-66','700000000000006'),
  ('P000007','Juliana Lima',       '1988-12-12','female','Brasilia','DF','777.777.777-77','700000000000007'),
  ('P000008','Rafael Almeida',     '1972-04-18','male',  'Formosa', 'GO','888.888.888-88','700000000000008'),
  ('P000009','Fernanda Rocha',     '1966-06-25','female','Brasilia','DF','999.999.999-99','700000000000009'),
  ('P000010','Lucas Martins',      '1990-08-08','male',  'Brasilia','DF','101.010.101-01','700000000000010'),
  ('P000011','Camila Ferreira',    '1948-02-14','female','Goiania', 'GO','121.212.121-21','700000000000011'),
  ('P000012','Gustavo Nunes',      '1983-10-05','male',  'Brasilia','DF','131.313.131-31','700000000000012'),
  ('P000013','Renata Dias',        '2000-05-19','female','Brasilia','DF','141.414.141-41','700000000000013'),
  ('P000014','Marcos Vieira',      '1958-03-03','male',  'Brasilia','DF','151.515.151-51','700000000000014'),
  ('P000015','Patricia Gomes',     '1975-11-11','female','Brasilia','DF','161.616.161-61','700000000000015');

-- ---------- encounters ----------
-- Um atendimento principal por paciente; diabéticos em Endocrinologia, hipertensos em Cardiologia.
INSERT INTO encounters (encounter_id, patient_id, start_date, end_date, encounter_type, department) VALUES
  ('ENC01','P000001','2024-02-10','2024-02-10','Ambulatorial','Endocrinologia'),
  ('ENC02','P000002','2024-03-05','2024-03-05','Retorno','Endocrinologia'),
  ('ENC03','P000003','2024-01-20','2024-01-22','Internacao','Endocrinologia'),
  ('ENC04','P000004','2024-04-11','2024-04-11','Ambulatorial','Endocrinologia'),
  ('ENC05','P000005','2024-05-30','2024-05-30','Ambulatorial','Endocrinologia'),
  ('ENC06','P000006','2024-02-28','2024-03-02','Internacao','Endocrinologia'),
  ('ENC07','P000007','2024-06-12','2024-06-12','Ambulatorial','Endocrinologia'),
  ('ENC08','P000008','2024-03-18','2024-03-18','Retorno','Endocrinologia'),
  ('ENC09','P000009','2024-04-25','2024-04-27','Internacao','Endocrinologia'),
  ('ENC10','P000010','2024-07-08','2024-07-08','Ambulatorial','Endocrinologia'),
  ('ENC11','P000011','2024-01-14','2024-01-16','Internacao','Endocrinologia'),
  ('ENC12','P000012','2024-08-05','2024-08-05','Ambulatorial','Endocrinologia'),
  ('ENC13','P000013','2024-05-19','2024-05-19','Ambulatorial','Cardiologia'),
  ('ENC14','P000014','2024-03-03','2024-03-03','Retorno','Cardiologia'),
  ('ENC15','P000015','2024-09-11','2024-09-11','Ambulatorial','Endocrinologia'),
  ('ENC16','P000001','2024-08-01','2024-08-01','Retorno','Cardiologia'),
  ('ENC17','P000007','2024-10-02','2024-10-02','Retorno','Cardiologia');

-- ---------- clinical_events: Condições ----------
INSERT INTO clinical_events (event_id, patient_id, encounter_id, event_type, code, description, event_date, value, unit) VALUES
  ('EVC01','P000001','ENC01','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-02-10',NULL,NULL),
  ('EVC02','P000002','ENC02','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-03-05',NULL,NULL),
  ('EVC03','P000003','ENC03','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-01-20',NULL,NULL),
  ('EVC04','P000004','ENC04','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-04-11',NULL,NULL),
  ('EVC05','P000005','ENC05','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-05-30',NULL,NULL),
  ('EVC06','P000006','ENC06','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-02-28',NULL,NULL),
  ('EVC07','P000007','ENC07','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-06-12',NULL,NULL),
  ('EVC08','P000008','ENC08','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-03-18',NULL,NULL),
  ('EVC09','P000009','ENC09','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-04-25',NULL,NULL),
  ('EVC10','P000010','ENC10','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-07-08',NULL,NULL),
  ('EVC11','P000011','ENC11','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-01-14',NULL,NULL),
  ('EVC12','P000012','ENC12','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-08-05',NULL,NULL),
  ('EVC13','P000015','ENC15','Condition','Diabetes',   'Diabetes Mellitus Tipo 2','2024-09-11',NULL,NULL),
  ('EVC14','P000001','ENC16','Condition','Hipertensao','Hipertensão Arterial',    '2024-08-01',NULL,NULL),
  ('EVC15','P000007','ENC17','Condition','Hipertensao','Hipertensão Arterial',    '2024-10-02',NULL,NULL),
  ('EVC16','P000013','ENC13','Condition','Hipertensao','Hipertensão Arterial',    '2024-05-19',NULL,NULL),
  ('EVC17','P000014','ENC14','Condition','Hipertensao','Hipertensão Arterial',    '2024-03-03',NULL,NULL);

-- ---------- clinical_events: Observações (exames) ----------
-- HbA1c dos 13 diabéticos (unidade %); valores escolhidos p/ média ~7.76 e mediana 7.6.
INSERT INTO clinical_events (event_id, patient_id, encounter_id, event_type, code, description, event_date, value, unit) VALUES
  ('EVO01','P000001','ENC01','Observation','HbA1c','Hemoglobina glicada','2024-02-10',8.1,'%'),
  ('EVO02','P000002','ENC02','Observation','HbA1c','Hemoglobina glicada','2024-03-05',7.2,'%'),
  ('EVO03','P000003','ENC03','Observation','HbA1c','Hemoglobina glicada','2024-01-20',9.0,'%'),
  ('EVO04','P000004','ENC04','Observation','HbA1c','Hemoglobina glicada','2024-04-11',6.8,'%'),
  ('EVO05','P000005','ENC05','Observation','HbA1c','Hemoglobina glicada','2024-05-30',7.5,'%'),
  ('EVO06','P000006','ENC06','Observation','HbA1c','Hemoglobina glicada','2024-02-28',8.8,'%'),
  ('EVO07','P000007','ENC07','Observation','HbA1c','Hemoglobina glicada','2024-06-12',7.0,'%'),
  ('EVO08','P000008','ENC08','Observation','HbA1c','Hemoglobina glicada','2024-03-18',6.5,'%'),
  ('EVO09','P000009','ENC09','Observation','HbA1c','Hemoglobina glicada','2024-04-25',9.4,'%'),
  ('EVO10','P000010','ENC10','Observation','HbA1c','Hemoglobina glicada','2024-07-08',7.8,'%'),
  ('EVO11','P000011','ENC11','Observation','HbA1c','Hemoglobina glicada','2024-01-14',8.3,'%'),
  ('EVO12','P000012','ENC12','Observation','HbA1c','Hemoglobina glicada','2024-08-05',6.9,'%'),
  ('EVO13','P000015','ENC15','Observation','HbA1c','Hemoglobina glicada','2024-09-11',7.6,'%'),
  -- Glicemia de jejum para alguns pacientes (unidade mg/dL).
  ('EVO14','P000001','ENC01','Observation','Glicemia','Glicemia de jejum','2024-02-10',182,'mg/dL'),
  ('EVO15','P000003','ENC03','Observation','Glicemia','Glicemia de jejum','2024-01-20',210,'mg/dL'),
  ('EVO16','P000006','ENC06','Observation','Glicemia','Glicemia de jejum','2024-02-28',198,'mg/dL'),
  ('EVO17','P000009','ENC09','Observation','Glicemia','Glicemia de jejum','2024-04-25',225,'mg/dL'),
  ('EVO18','P000011','ENC11','Observation','Glicemia','Glicemia de jejum','2024-01-14',190,'mg/dL');

-- ---------- clinical_events: Medicações ----------
-- Metformina (10 na coorte), Insulina (4 na coorte), Losartana (2 na coorte + 2 fora).
INSERT INTO clinical_events (event_id, patient_id, encounter_id, event_type, code, description, event_date, value, unit) VALUES
  ('EVM01','P000001','ENC01','Medication','Metformina','Metformina 850 mg','2024-02-10',850,'mg'),
  ('EVM02','P000002','ENC02','Medication','Metformina','Metformina 850 mg','2024-03-05',850,'mg'),
  ('EVM03','P000003','ENC03','Medication','Metformina','Metformina 850 mg','2024-01-20',850,'mg'),
  ('EVM04','P000004','ENC04','Medication','Metformina','Metformina 850 mg','2024-04-11',850,'mg'),
  ('EVM05','P000006','ENC06','Medication','Metformina','Metformina 850 mg','2024-02-28',850,'mg'),
  ('EVM06','P000008','ENC08','Medication','Metformina','Metformina 850 mg','2024-03-18',850,'mg'),
  ('EVM07','P000009','ENC09','Medication','Metformina','Metformina 850 mg','2024-04-25',850,'mg'),
  ('EVM08','P000011','ENC11','Medication','Metformina','Metformina 850 mg','2024-01-14',850,'mg'),
  ('EVM09','P000012','ENC12','Medication','Metformina','Metformina 850 mg','2024-08-05',850,'mg'),
  ('EVM10','P000015','ENC15','Medication','Metformina','Metformina 850 mg','2024-09-11',850,'mg'),
  ('EVM11','P000003','ENC03','Medication','Insulina',  'Insulina NPH 10 UI',  '2024-01-20',10,'UI'),
  ('EVM12','P000006','ENC06','Medication','Insulina',  'Insulina NPH 10 UI',  '2024-02-28',10,'UI'),
  ('EVM13','P000009','ENC09','Medication','Insulina',  'Insulina NPH 10 UI',  '2024-04-25',10,'UI'),
  ('EVM14','P000011','ENC11','Medication','Insulina',  'Insulina NPH 10 UI',  '2024-01-14',10,'UI'),
  ('EVM15','P000001','ENC16','Medication','Losartana', 'Losartana 50 mg',     '2024-08-01',50,'mg'),
  ('EVM16','P000007','ENC17','Medication','Losartana', 'Losartana 50 mg',     '2024-10-02',50,'mg'),
  ('EVM17','P000013','ENC13','Medication','Losartana', 'Losartana 50 mg',     '2024-05-19',50,'mg'),
  ('EVM18','P000014','ENC14','Medication','Losartana', 'Losartana 50 mg',     '2024-03-03',50,'mg');

-- ---------- user_patient_assignments ----------
-- med.cardoso: médico de P000001..P000008.
-- est.souza: estagiário (supervisionado por med.cardoso) de P000001..P000003.
-- med.almeida: médico de P000009..P000015.
INSERT INTO user_patient_assignments (assignment_id, caregiver_username, patient_id, assignment_type, supervisor_username, status) VALUES
  ('A01','med.cardoso','P000001','medico',    NULL,          'ativo'),
  ('A02','med.cardoso','P000002','medico',    NULL,          'ativo'),
  ('A03','med.cardoso','P000003','medico',    NULL,          'ativo'),
  ('A04','med.cardoso','P000004','medico',    NULL,          'ativo'),
  ('A05','med.cardoso','P000005','medico',    NULL,          'ativo'),
  ('A06','med.cardoso','P000006','medico',    NULL,          'ativo'),
  ('A07','med.cardoso','P000007','medico',    NULL,          'ativo'),
  ('A08','med.cardoso','P000008','medico',    NULL,          'ativo'),
  ('A09','est.souza',  'P000001','estagiario','med.cardoso', 'ativo'),
  ('A10','est.souza',  'P000002','estagiario','med.cardoso', 'ativo'),
  ('A11','est.souza',  'P000003','estagiario','med.cardoso', 'ativo'),
  ('A12','med.almeida','P000009','medico',    NULL,          'ativo'),
  ('A13','med.almeida','P000010','medico',    NULL,          'ativo'),
  ('A14','med.almeida','P000011','medico',    NULL,          'ativo'),
  ('A15','med.almeida','P000012','medico',    NULL,          'ativo'),
  ('A16','med.almeida','P000013','medico',    NULL,          'ativo'),
  ('A17','med.almeida','P000014','medico',    NULL,          'ativo'),
  ('A18','med.almeida','P000015','medico',    NULL,          'ativo');

-- ---------- projects ----------
INSERT INTO projects (project_id, title, researcher_username, condition_code, status, valid_until) VALUES
  ('PRJ01','Coorte Diabetes Tipo 2',        'pesq.lima','Diabetes',   'Aprovado','2027-12-31'),
  ('PRJ02','Coorte Hipertensão Resistente', 'pesq.lima','Hipertensao','Expirado','2024-01-01'),
  ('PRJ03','Coorte Obesidade',              'pesq.melo','Obesidade',  'Aprovado','2027-06-30');
