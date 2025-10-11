-- Add users table
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add unique constraint to enforce one account per currency per user
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency" UNIQUE ("owner", "currency");
 ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- Add indexes for users table
CREATE INDEX ON "users" ("username");
CREATE INDEX ON "users" ("email");

