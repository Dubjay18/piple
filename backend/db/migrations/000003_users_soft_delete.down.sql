DROP INDEX IF EXISTS "users_email_active_key";
ALTER TABLE "users" ADD CONSTRAINT "users_email_key" UNIQUE ("email");
DROP INDEX IF EXISTS "users_deleted_at_idx";
ALTER TABLE "users" DROP COLUMN IF EXISTS "deleted_at";