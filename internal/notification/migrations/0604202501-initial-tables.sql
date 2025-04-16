-- *****************************
-- ********** SCHEMA ***********
-- *****************************

CREATE SCHEMA IF NOT EXISTS notification;

-- *****************************
-- ****** SUBSCRIPTIONS ********
-- *****************************

DROP TABLE IF EXISTS notification.subscriptions;

CREATE TABLE IF NOT EXISTS notification.subscriptions
(
    project_id uuid NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    language character varying(4) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL
);

ALTER TABLE IF EXISTS notification.subscriptions OWNER to admin;

-- ******************************
-- ****** EMAIL ACCOUNTS ********
-- ******************************

DROP TABLE IF EXISTS notification.email_accounts;

CREATE TABLE IF NOT EXISTS notification.email_accounts
(
    id uuid NOT NULL,
    project_id uuid NOT NULL,
    email character varying(128) COLLATE pg_catalog."default" NOT NULL,
    display_name character varying(64) COLLATE pg_catalog."default" NOT NULL,
    host text COLLATE pg_catalog."default",
    port smallint,
    username text COLLATE pg_catalog."default",
    password text COLLATE pg_catalog."default",
    enable_ssl boolean NOT NULL DEFAULT false,
    email_authentication_method_type_id smallint NOT NULL,
    client_id text COLLATE pg_catalog."default",
    client_secret text COLLATE pg_catalog."default",
    tenant_id text COLLATE pg_catalog."default",
    created_at timestamp without time zone NOT NULL,
    CONSTRAINT "PK_email_accounts" PRIMARY KEY (id),
    CONSTRAINT "UX_email_accounts_project_id_email" UNIQUE (project_id, email)
);

ALTER TABLE IF EXISTS notification.email_accounts OWNER to admin;

-- *******************************
-- ****** EMAIL TEMPLATES ********
-- *******************************

DROP TABLE IF EXISTS notification.email_templates;

CREATE TABLE IF NOT EXISTS notification.email_templates
(
    email_account_id uuid NOT NULL,
    name character varying(128) COLLATE pg_catalog."default" NOT NULL,
    language character varying(8) COLLATE pg_catalog."default" NOT NULL,
    subject character varying(128) COLLATE pg_catalog."default" NOT NULL,
    body text COLLATE pg_catalog."default" NOT NULL,
    bcc_email_addresses text COLLATE pg_catalog."default",
    allow_direct_reply boolean NOT NULL DEFAULT false,
    CONSTRAINT "UX_email_templates_email_account_id_name_language" UNIQUE (email_account_id, name, language),
    CONSTRAINT "FK_email_templates_email_account_id" FOREIGN KEY (email_account_id)
        REFERENCES public.email_accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS notification.email_templates OWNER to admin;

-- *****************************
-- ****** QUEUED EMAILS ********
-- *****************************

DROP TABLE IF EXISTS notification.queued_emails;

CREATE TABLE IF NOT EXISTS notification.queued_emails
(
    email_account_id uuid NOT NULL,
    "to" character varying(128) COLLATE pg_catalog."default" NOT NULL,
    reply_to character varying(128) COLLATE pg_catalog."default",
    cc text COLLATE pg_catalog."default",
    bcc text COLLATE pg_catalog."default",
    subject character varying(128) COLLATE pg_catalog."default" NOT NULL,
    body text COLLATE pg_catalog."default" NOT NULL,
    sent_at timestamp without time zone,
    sent_tries smallint NOT NULL,
    CONSTRAINT "FK_queued_emails_email_account_id" FOREIGN KEY (email_account_id)
        REFERENCES public.email_accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS notification.queued_emails OWNER to admin;