ALTER TABLE "users" ADD COLUMN "deleted_at" timestamp;

CREATE INDEX "users_deleted_at_idx" ON "users" ("deleted_at");

-- Soft delete keeps the row around, so the plain UNIQUE(email) has to go —
-- otherwise re-provisioning someone at the same email after their old
-- account was deleted would fail. Partial index enforces uniqueness only
-- among live accounts.
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_email_key";
CREATE UNIQUE INDEX "users_email_active_key" ON "users" ("email") WHERE "deleted_at" IS NULL;