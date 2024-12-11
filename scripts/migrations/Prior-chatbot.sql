DROP TABLE IF EXISTS PC_TI_USER;

CREATE TABLE PC_TI_USER
(
    USER_ID    SERIAL PRIMARY KEY,
    EMAIL      VARCHAR(255) NOT NULL,
    PASSWORD   VARCHAR(255) NOT NULL,
    MOBILE     VARCHAR(255),
    DOB        DATE,
    SEX        VARCHAR(10),
    CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UPDATED_AT TIMESTAMP,
    IS_DELETED CHAR DEFAULT 'N'::BPCHAR NOT NULL
);

INSERT INTO PC_TI_USER (EMAIL, PASSWORD, MOBILE, DOB, SEX)
VALUES ('admin@prior-chatbot','$2a$12$dyyRdvdlUN1b7p9kibwlu.sWgPj5C2o5fzdWuzGyfyPHJPOb5.jyS',
        '0812345678','1990-01-01','M');
-- adminprior-chatbot@PassW0rd = password