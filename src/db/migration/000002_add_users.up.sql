CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "deslies" ADD "owner" varchar(255) NOT NULL;

CREATE INDEX ON "deslies" ("owner");

CREATE INDEX ON "deslies" ("desly");

CREATE INDEX ON "deslies" ("redirect");

ALTER TABLE "deslies" ADD CONSTRAINT "owner_redirect_key" UNIQUE ("owner", "redirect");

ALTER TABLE "deslies" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");