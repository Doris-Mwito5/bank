-- Drop foreign key constraints first
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- Then drop the unique constraint if it exists
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency";

-- Finally drop the users table
DROP TABLE IF EXISTS "users";