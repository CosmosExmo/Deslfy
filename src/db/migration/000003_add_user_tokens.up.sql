CREATE TABLE "user_tokens" (
  "id" SERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "token" varchar NOT NULL,
  "expire_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "user_tokens" ("owner");

CREATE INDEX ON "user_tokens" ("token");

CREATE INDEX ON "user_tokens" ("expire_at");

CREATE INDEX ON "user_tokens" ("owner", "token");

CREATE INDEX ON "user_tokens" ("owner", "token", "expire_at");

ALTER TABLE "user_tokens" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

