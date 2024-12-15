DROP TABLE IF EXISTS "verify_emails" CASCADE;

ALTER TABLE IF EXISTS "users" DROP COLUMN "is_email_verified"
