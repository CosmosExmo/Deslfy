CREATE TABLE "deslies" (
  "id" serial PRIMARY KEY,
  "redirect" varchar NOT NULL,
  "desly" varchar NOT NULL UNIQUE,
  "clicked" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "deslies" ("desly");