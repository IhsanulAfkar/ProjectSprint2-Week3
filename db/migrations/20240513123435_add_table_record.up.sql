CREATE TABLE IF NOT EXISTS "record" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "identityNumber" bigint NOT NULL,
    "nip" bigint NOT NULL,
    "symptoms" TEXT NOT NULL,
    "medications" TEXT NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "record" ADD FOREIGN KEY ("identityNumber") REFERENCES "patient" ("identityNumber");