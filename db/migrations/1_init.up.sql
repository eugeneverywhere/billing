SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;
CREATE TABLE accounts (
  id                    BIGINT AUTO_INCREMENT   PRIMARY KEY,

  external_id           VARCHAR(20)  NOT NULL UNIQUE KEY ,
  balance               FLOAT            NOT NULL,

  status             SMALLINT(6)          default '10' not null,
  created_at         TIMESTAMP            NOT NULL DEFAULT NOW(),
  updated_at         TIMESTAMP            NOT NULL DEFAULT NOW() ON UPDATE NOW()
)
  CHARACTER SET utf8
  COLLATE utf8_unicode_ci
  ENGINE = InnoDB;
