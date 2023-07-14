-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-07-13T17:49:35.605Z

CREATE TABLE "deslies" (
  "id" SERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "redirect" varchar NOT NULL,
  "desly" varchar UNIQUE NOT NULL,
  "clicked" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_email_verified" boolean NOT NULL DEFAULT false,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

CREATE TABLE "user_tokens" (
  "id" SERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "token" varchar NOT NULL,
  "expire_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "deslies" ("owner");

CREATE INDEX ON "deslies" ("desly");

CREATE INDEX ON "deslies" ("redirect");

CREATE INDEX ON "deslies" ("owner", "desly");

CREATE UNIQUE INDEX ON "deslies" ("owner", "redirect");

CREATE INDEX ON "user_tokens" ("owner");

CREATE INDEX ON "user_tokens" ("token");

CREATE INDEX ON "user_tokens" ("expire_at");

CREATE INDEX ON "user_tokens" ("owner", "token");

CREATE UNIQUE INDEX ON "user_tokens" ("owner", "token", "expire_at");

CREATE INDEX ON "sessions" ("id");

CREATE INDEX ON "sessions" ("username");

ALTER TABLE "deslies" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "user_tokens" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
