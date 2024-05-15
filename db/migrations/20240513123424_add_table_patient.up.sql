DROP TYPE IF EXISTS gender;
CREATE TYPE gender AS ENUM (
    'male',
    'female'
);

CREATE TABLE IF NOT EXISTS "patient" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "identityNumber" bigint UNIQUE NOT NULL,
    "name" varchar(50) NOT NULL,
    "phoneNumber" varchar(20) NOT NULL UNIQUE,
    "birthDate" DATE NOT NULL,
    "gender" gender NOT NULL,
    "identityCardScanImg" varchar(255) NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT (now())
);