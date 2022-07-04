-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "email" TEXT NOT NULl,
    "password" TEXT NOT NULl,
    "created_at" timestamptz DEFAULT now(),
    "updated_at" timestamptz DEFAULT now()
);

CREATE UNIQUE INDEX "uidx_users_email" ON "users" USING BTREE (LOWER("email"));

-- +migrate Down
DROP TABLE IF EXISTS "users";

DROP EXTENSION IF EXISTS "uuid-ossp";
