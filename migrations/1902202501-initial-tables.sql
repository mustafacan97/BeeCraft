-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id                      uuid PRIMARY KEY,    
    first_name              character varying(128) COLLATE pg_catalog."default",
    last_name               character varying(128) COLLATE pg_catalog."default",
    email                   character varying(128) COLLATE pg_catalog."default" NOT NULL,
    email_validated         boolean NOT NULL,
    phone                   character varying(16) COLLATE pg_catalog."default",
    phone_validated         boolean NOT NULL,
    gender_type_id          integer,
    birth_date              timestamp without time zone,
    password_hash           text COLLATE pg_catalog."default" NOT NULL,
    last_password_change_at timestamp without time zone,
    failed_login_attempts   integer NOT NULL,
    cannot_login_until_at   timestamp without time zone,
    refresh_token           text COLLATE pg_catalog."default",
    refresh_token_expire_at timestamp without time zone,
    last_ip_address         character varying(16) COLLATE pg_catalog."default",
    last_login_at           timestamp without time zone,
    is_system_user          boolean NOT NULL,
    admin_comment           character varying(512) COLLATE pg_catalog."default",
    created_at              timestamp without time zone NOT NULL,
    active                  boolean NOT NULL,
    deleted                 boolean NOT NULL
);

ALTER TABLE IF EXISTS public.users
    OWNER to admin;

-- Table: public.project

-- DROP TABLE IF EXISTS public.project;

CREATE TABLE IF NOT EXISTS public.projects
(
    id          uuid PRIMARY KEY,
    owner_id    uuid NOT NULL,
    name        character varying(128) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    url         text COLLATE pg_catalog."default",
    modules     text COLLATE pg_catalog."default",
    created_at  timestamp without time zone NOT NULL,
    active      boolean NOT NULL,
    deleted     boolean NOT NULL,

    CONSTRAINT "FK_projects_owner_id" FOREIGN KEY (owner_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.projects
    OWNER to admin;

-- Table: public.api_keys

-- DROP TABLE IF EXISTS public.api_keys;

CREATE TABLE IF NOT EXISTS public.api_keys
(
    project_id     uuid NOT NULL,
    hashed_api_key text COLLATE pg_catalog."default" NOT NULL,
    expired_at     timestamp without time zone,
    created_at     timestamp without time zone NOT NULL,
    active         boolean NOT NULL,
    deleted        boolean NOT NULL,

    CONSTRAINT "FK_api_keys_project_id" FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.api_keys
    OWNER to admin;

-- Table: public.otp_code

-- DROP TABLE IF EXISTS public.otp_code;

CREATE TABLE IF NOT EXISTS public.otp_codes
(
    id          SERIAL PRIMARY KEY,
    hashed_code text COLLATE pg_catalog."default" NOT NULL,
    user_id     uuid NOT NULL,
    type_id     smallint NOT NULL,
    expired_at  timestamp without time zone NOT NULL,
    created_at  timestamp without time zone NOT NULL,
    
    CONSTRAINT "FK_otp_codes_user_id" FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.otp_codes
    OWNER to admin;

-- Table: public.permissions

-- DROP TABLE IF EXISTS public.permissions;

CREATE TABLE IF NOT EXISTS public.permissions
(
    id         SERIAL PRIMARY KEY,
    name       character varying(256) COLLATE pg_catalog."default" NOT NULL,
    service_id integer
);

ALTER TABLE IF EXISTS public.permissions
    OWNER to admin;

-- Index: IX_permission_name_service_id

-- DROP INDEX IF EXISTS public."IX_permission_name_service_id";

CREATE UNIQUE INDEX IF NOT EXISTS "IX_permission_name_service_id"
    ON public.permissions USING btree
    (name COLLATE pg_catalog."default" ASC NULLS LAST, service_id ASC NULLS LAST);

-- Table: public.roles

-- DROP TABLE IF EXISTS public.roles;

CREATE TABLE IF NOT EXISTS public.roles
(
    id         SERIAL PRIMARY KEY,
    name       character varying(64) COLLATE pg_catalog."default" NOT NULL,
    project_id uuid,

    CONSTRAINT "FK_roles_project_id" FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.roles
    OWNER to admin;

-- Index: IX_role_name

-- DROP INDEX IF EXISTS public."IX_role_name";

CREATE UNIQUE INDEX IF NOT EXISTS "IX_role_name"
    ON public.roles USING btree
        (name COLLATE pg_catalog."default" ASC NULLS LAST);

-- Table: public.outbox_messages

-- DROP TABLE IF EXISTS public.outbox_messages;

CREATE TABLE IF NOT EXISTS public.outbox_messages
(
    id           uuid PRIMARY KEY,
    message_type character varying(256) COLLATE pg_catalog."default" NOT NULL,
    payload      text COLLATE pg_catalog."default" NOT NULL,
    error        text COLLATE pg_catalog."default" NULL,
    created_at   timestamp without time zone NOT NULL,
    processed_at timestamp without time zone
);

ALTER TABLE IF EXISTS public.outbox_message
    OWNER to admin;

-- Table: public.permission_role_mappings

-- DROP TABLE IF EXISTS public.permission_role_mappings;

CREATE TABLE IF NOT EXISTS public.permission_role_mapping
(
    permission_id INTEGER NOT NULL,
    role_id       INTEGER NOT NULL,

    CONSTRAINT "PK_permission_role_mappings" PRIMARY KEY (permission_id, role_id),
    CONSTRAINT "FK_permission_role_mappings_permission_id" FOREIGN KEY (permission_id)
        REFERENCES public.permissions (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT "FK_permission_role_mappings_role_id" FOREIGN KEY (role_id)
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.permission_role_mapping
    OWNER to admin;

-- Table: public.user_role_mappings

-- DROP TABLE IF EXISTS public.user_role_mappings;

CREATE TABLE IF NOT EXISTS public.user_role_mappings
(
    user_id uuid NOT NULL,
    role_id INTEGER NOT NULL,

    CONSTRAINT "PK_user_role_mappings" PRIMARY KEY (user_id, role_id),
    CONSTRAINT "FK_user_role_mapping_role_id" FOREIGN KEY (role_id)
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT "FK_user_role_mappings_user_id" FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.user_role_mapping
    OWNER to admin;

-- Table: public.user_project_mappings

-- DROP TABLE IF EXISTS public.user_project_mappings;

CREATE TABLE IF NOT EXISTS public.project_members
(
    user_id    uuid NOT NULL,
    project_id uuid NOT NULL,

    CONSTRAINT "PK_project_members" PRIMARY KEY (user_id, project_id),
    CONSTRAINT "FK_project_members_user_id" FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT "FK_project_members_project_id" FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.user_role_mapping
    OWNER to admin;