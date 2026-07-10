-- Dados de exemplo para desenvolvimento/teste local, usando as MESMAS convenções do
-- banco real (event_type/códigos/departamentos em MAIÚSCULO/inglês, value como texto,
-- assignment_type ATTENDING/TRAINEE, coluna active booleana, target_condition_code).
-- Coorte "DIABETES": 13 pacientes; sexo 7F/6M; faixas 18-39=3, 40-59=6, 60+=4;
-- HbA1c média ~7.76 / mediana 7.6; METFORMIN=10, INSULIN=4, LOSARTAN=2;
-- ENDOCRINOLOGY=13, CARDIOLOGY=2.

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
INSERT INTO encounters (encounter_id, patient_id, start_date, end_date, encounter_type, department) VALUES
  ('ENC01','P000001','2024-02-10','2024-02-10','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC02','P000002','2024-03-05','2024-03-05','FOLLOW_UP','ENDOCRINOLOGY'),
  ('ENC03','P000003','2024-01-20','2024-01-22','INPATIENT','ENDOCRINOLOGY'),
  ('ENC04','P000004','2024-04-11','2024-04-11','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC05','P000005','2024-05-30','2024-05-30','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC06','P000006','2024-02-28','2024-03-02','INPATIENT','ENDOCRINOLOGY'),
  ('ENC07','P000007','2024-06-12','2024-06-12','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC08','P000008','2024-03-18','2024-03-18','FOLLOW_UP','ENDOCRINOLOGY'),
  ('ENC09','P000009','2024-04-25','2024-04-27','INPATIENT','ENDOCRINOLOGY'),
  ('ENC10','P000010','2024-07-08','2024-07-08','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC11','P000011','2024-01-14','2024-01-16','INPATIENT','ENDOCRINOLOGY'),
  ('ENC12','P000012','2024-08-05','2024-08-05','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC13','P000013','2024-05-19','2024-05-19','AMBULATORIAL','CARDIOLOGY'),
  ('ENC14','P000014','2024-03-03','2024-03-03','FOLLOW_UP','CARDIOLOGY'),
  ('ENC15','P000015','2024-09-11','2024-09-11','AMBULATORIAL','ENDOCRINOLOGY'),
  ('ENC16','P000001','2024-08-01','2024-08-01','FOLLOW_UP','CARDIOLOGY'),
  ('ENC17','P000007','2024-10-02','2024-10-02','FOLLOW_UP','CARDIOLOGY');

-- ---------- clinical_events: CONDITION ----------
INSERT INTO clinical_events (event_id, patient_id, encounter_id, event_type, code, description, value, unit, event_date) VALUES
  ('EVC01','P000001','ENC01','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-02-10'),
  ('EVC02','P000002','ENC02','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-03-05'),
  ('EVC03','P000003','ENC03','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-01-20'),
  ('EVC04','P000004','ENC04','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-04-11'),
  ('EVC05','P000005','ENC05','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-05-30'),
  ('EVC06','P000006','ENC06','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-02-28'),
  ('EVC07','P000007','ENC07','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-06-12'),
  ('EVC08','P000008','ENC08','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-03-18'),
  ('EVC09','P000009','ENC09','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-04-25'),
  ('EVC10','P000010','ENC10','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-07-08'),
  ('EVC11','P000011','ENC11','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-01-14'),
  ('EVC12','P000012','ENC12','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-08-05'),
  ('EVC13','P000015','ENC15','CONDITION','DIABETES',    'Diabetes Mellitus Tipo 2',NULL,NULL,'2024-09-11'),
  ('EVC14','P000001','ENC16','CONDITION','HYPERTENSION','Hipertensão Arterial',    NULL,NULL,'2024-08-01'),
  ('EVC15','P000007','ENC17','CONDITION','HYPERTENSION','Hipertensão Arterial',    NULL,NULL,'2024-10-02'),
  ('EVC16','P000013','ENC13','CONDITION','HYPERTENSION','Hipertensão Arterial',    NULL,NULL,'2024-05-19'),
  ('EVC17','P000014','ENC14','CONDITION','HYPERTENSION','Hipertensão Arterial',    NULL,NULL,'2024-03-03');

-- ---------- clinical_events: OBSERVATION ----------
INSERT INTO clinical_events (event_id, patient_id, encounter_id, event_type, code, description, value, unit, event_date) VALUES
  ('EVO01','P000001','ENC01','OBSERVATION','HBA1C','Hemoglobina glicada','8.1','%','2024-02-10'),
  ('EVO02','P000002','ENC02','OBSERVATION','HBA1C','Hemoglobina glicada','7.2','%','2024-03-05'),
  ('EVO03','P000003','ENC03','OBSERVATION','HBA1C','Hemoglobina glicada','9.0','%','2024-01-20'),
  ('EVO04','P000004','ENC04','OBSERVATION','HBA1C','Hemoglobina glicada','6.8','%','2024-04-11'),
  ('EVO05','P000005','ENC05','OBSERVATION','HBA1C','Hemoglobina glicada','7.5','%','2024-05-30'),
  ('EVO06','P000006','ENC06','OBSERVATION','HBA1C','Hemoglobina glicada','8.8','%','2024-02-28'),
  ('EVO07','P000007','ENC07','OBSERVATION','HBA1C','Hemoglobina glicada','7.0','%','2024-06-12'),
  ('EVO08','P000008','ENC08','OBSERVATION','HBA1C','Hemoglobina glicada','6.5','%','2024-03-18'),
  ('EVO09','P000009','ENC09','OBSERVATION','HBA1C','Hemoglobina glicada','9.4','%','2024-04-25'),
  ('EVO10','P000010','ENC10','OBSERVATION','HBA1C','Hemoglobina glicada','7.8','%','2024-07-08'),
  ('EVO11','P000011','ENC11','OBSERVATION','HBA1C','Hemoglobina glicada','8.3','%','2024-01-14'),
  ('EVO12','P000012','ENC12','OBSERVATION','HBA1C','Hemoglobina glicada','6.9','%','2024-08-05'),
  ('EVO13','P000015','ENC15','OBSERVATION','HBA1C','Hemoglobina glicada','7.6','%','2024-09-11'),
  ('EVO14','P000001','ENC01','OBSERVATION','GLUCOSE','Glicemia de jejum','182','mg/dL','2024-02-10'),
  ('EVO15','P000003','ENC03','OBSERVATION','GLUCOSE','Glicemia de jejum','210','mg/dL','2024-01-20'),
  ('EVO16','P000006','ENC06','OBSERVATION','GLUCOSE','Glicemia de jejum','198','mg/dL','2024-02-28'),
  ('EVO17','P000009','ENC09','OBSERVATION','GLUCOSE','Glicemia de jejum','225','mg/dL','2024-04-25'),
  ('EVO18','P000011','ENC11','OBSERVATION','GLUCOSE','Glicemia de jejum','190','mg/dL','2024-01-14');

-- ---------- clinical_events: MEDICATION ----------
INSERT INTO clinical_events (event_id, patient_id, encounter_id, event_type, code, description, value, unit, event_date) VALUES
  ('EVM01','P000001','ENC01','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-02-10'),
  ('EVM02','P000002','ENC02','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-03-05'),
  ('EVM03','P000003','ENC03','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-01-20'),
  ('EVM04','P000004','ENC04','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-04-11'),
  ('EVM05','P000006','ENC06','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-02-28'),
  ('EVM06','P000008','ENC08','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-03-18'),
  ('EVM07','P000009','ENC09','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-04-25'),
  ('EVM08','P000011','ENC11','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-01-14'),
  ('EVM09','P000012','ENC12','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-08-05'),
  ('EVM10','P000015','ENC15','MEDICATION','METFORMIN','Metformina 850 mg','850','mg','2024-09-11'),
  ('EVM11','P000003','ENC03','MEDICATION','INSULIN',  'Insulina NPH 10 UI',  '10','UI','2024-01-20'),
  ('EVM12','P000006','ENC06','MEDICATION','INSULIN',  'Insulina NPH 10 UI',  '10','UI','2024-02-28'),
  ('EVM13','P000009','ENC09','MEDICATION','INSULIN',  'Insulina NPH 10 UI',  '10','UI','2024-04-25'),
  ('EVM14','P000011','ENC11','MEDICATION','INSULIN',  'Insulina NPH 10 UI',  '10','UI','2024-01-14'),
  ('EVM15','P000001','ENC16','MEDICATION','LOSARTAN', 'Losartana 50 mg',     '50','mg','2024-08-01'),
  ('EVM16','P000007','ENC17','MEDICATION','LOSARTAN', 'Losartana 50 mg',     '50','mg','2024-10-02'),
  ('EVM17','P000013','ENC13','MEDICATION','LOSARTAN', 'Losartana 50 mg',     '50','mg','2024-05-19'),
  ('EVM18','P000014','ENC14','MEDICATION','LOSARTAN', 'Losartana 50 mg',     '50','mg','2024-03-03');

-- ---------- user_patient_assignments ----------
-- med.cardoso: ATTENDING de P000001..P000008.
-- est.souza: TRAINEE (supervisionado por med.cardoso) de P000001..P000003.
-- med.almeida: ATTENDING de P000009..P000015.
INSERT INTO user_patient_assignments (assignment_id, username, patient_id, assignment_type, supervisor_username, active) VALUES
  ('A01','med.cardoso','P000001','ATTENDING',NULL,          TRUE),
  ('A02','med.cardoso','P000002','ATTENDING',NULL,          TRUE),
  ('A03','med.cardoso','P000003','ATTENDING',NULL,          TRUE),
  ('A04','med.cardoso','P000004','ATTENDING',NULL,          TRUE),
  ('A05','med.cardoso','P000005','ATTENDING',NULL,          TRUE),
  ('A06','med.cardoso','P000006','ATTENDING',NULL,          TRUE),
  ('A07','med.cardoso','P000007','ATTENDING',NULL,          TRUE),
  ('A08','med.cardoso','P000008','ATTENDING',NULL,          TRUE),
  ('A09','est.souza',  'P000001','TRAINEE',  'med.cardoso', TRUE),
  ('A10','est.souza',  'P000002','TRAINEE',  'med.cardoso', TRUE),
  ('A11','est.souza',  'P000003','TRAINEE',  'med.cardoso', TRUE),
  ('A12','med.almeida','P000009','ATTENDING',NULL,          TRUE),
  ('A13','med.almeida','P000010','ATTENDING',NULL,          TRUE),
  ('A14','med.almeida','P000011','ATTENDING',NULL,          TRUE),
  ('A15','med.almeida','P000012','ATTENDING',NULL,          TRUE),
  ('A16','med.almeida','P000013','ATTENDING',NULL,          TRUE),
  ('A17','med.almeida','P000014','ATTENDING',NULL,          TRUE),
  ('A18','med.almeida','P000015','ATTENDING',NULL,          TRUE);

-- ---------- projects ----------
INSERT INTO projects (project_id, title, researcher_username, target_condition_code, status, valid_until) VALUES
  ('PRJ01','Coorte Diabetes Tipo 2',        'pesq.lima','DIABETES',    'APPROVED','2027-12-31'),
  ('PRJ02','Coorte Hipertensão Resistente', 'pesq.lima','HYPERTENSION','EXPIRED', '2024-01-01'),
  ('PRJ03','Coorte Obesidade',              'pesq.melo','OBESITY',     'APPROVED','2027-06-30');
