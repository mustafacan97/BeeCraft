DO $$
    DECLARE
        project_id UUID := 'd3a99d5e-3c8c-4b39-84e2-a814da4db011';
        email_account_id UUID := '44e7ac3f-d914-4890-8e1a-91713c375219';
    BEGIN    
    
        INSERT INTO notification.email_accounts (id, project_id, email, display_name, host, port, enable_ssl, type_id, client_id, tenant_id, client_secret, created_at)
        VALUES 
            (email_account_id, project_id, 'mustafa@yeyu.co', 'mustafa@yeyu.co', 'smtp.office365.com', 587, 3, '', '', '', CURRENT_TIMESTAMP);
        
        INSERT INTO notification.email_templates (email_account_id, name, language, subject, body, allow_direct_reply)
        VALUES
            (email_account_id, 'USER_EMAIL_VALIDATION', 'tr-TR', 'E-Posta Doğrulama', 'Yeyu Platform''a hoş geldiniz!<br />Hesabını aktif etmek için buraya <a href="%AccountActivationURL%">tıklayın</a>.<br />', FALSE),
            (email_account_id, 'USER_EMAIL_VALIDATION', 'en-EN', 'Email Validation', 'Welcome to Yeyu Platform!<br />To activate your account <a href="%AccountActivationURL%">click here</a>.<br />', FALSE);

    END $$;