CREATE TABLE "services" (
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "name" varchar UNIQUE NOT NULL,
    "description" varchar NOT NULL,
    "created_at" timestamp DEFAULT (now()),
    "updated_at" timestamp DEFAULT (now()),
    "deleted_at" timestamp
);

CREATE TABLE "service_versions" (
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "changelog" varchar UNIQUE NOT NULL,
    "version" varchar NOT NULL,
    "service_id" uuid NOT NULL,
    "created_at" timestamp DEFAULT (now()),
    "updated_at" timestamp DEFAULT (now()),
    "deleted_at" timestamp
);

CREATE INDEX ON "services" ("id");

CREATE INDEX ON "services" ("name");

CREATE INDEX ON "service_versions" ("id");

CREATE INDEX ON "service_versions" ("changelog");

CREATE INDEX ON "service_versions" ("id", "service_id");

ALTER TABLE "service_versions" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
