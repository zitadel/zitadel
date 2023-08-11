ALTER TABLE auth.users ADD COLUMN otp_sms_added BOOL DEFAULT false;
ALTER TABLE auth.users ADD COLUMN otp_email_added BOOL DEFAULT false;