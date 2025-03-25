DO $$
    BEGIN

        -- Insert into role table
        INSERT INTO roles (name, project_id) 
        VALUES ('ADMIN', null), ('PROJECT_OWNER', null), ('REGISTERED', null);

    END $$;
