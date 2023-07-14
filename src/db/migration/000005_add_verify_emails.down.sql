ALTER TABLE IF EXISTS "verify_emails" DROP CONSTRAINT IF EXISTS "verify_emails_username_fkey";

DROP TABLE IF EXISTS "verify_emails";

ALTER TABLE "users" DROP COLUMN "is_email_verified";