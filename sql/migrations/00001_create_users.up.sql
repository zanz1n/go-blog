CREATE TYPE "UserRole" AS ENUM ('ADMIN', 'PUBLISHER', 'COMMON');

CREATE TABLE "users" (
    "id" uuid PRIMARY KEY ,
    "created_at" timestamp(3) NOT NULL DEFAULT current_timestamp,
    "updated_at" timestamp(3) NOT NULL DEFAULT current_timestamp,
    "email" varchar(64) NOT NULL,
    "username" varchar(42) NOT NULL,
    "password" varchar(60) NOT NULL,
    "role" "UserRole" NOT NULL DEFAULT 'COMMON'
);

CREATE UNIQUE INDEX "users_email_idx" ON "users"("email");
