CREATE TABLE IF NOT EXISTS "image" (
    "id" uuid UNIQUE NOT NULL DEFAULT (gen_random_uuid()) PRIMARY KEY,
    "path" varchar(255) NOT NULL,
    "createdAt" timestamp NOT NULL DEFAULT (now())
);