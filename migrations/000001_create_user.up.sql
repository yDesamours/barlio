CREATE TABLE IF NOT EXISTS "users"(
    "id" SERIAL PRIMARY KEY,
    "username" CITEXT UNIQUE NOT NULL,
    "email" TEXT UNIQUE NOT NULL,
    "lastname" TEXT,
    "firstname" TEXT,
    "isverified" boolean NOT NULL DEFAULT false,
    "joined_at" TIMESTAMP NOT NULL,
    "birthdate" TIMESTAMP,
    "bio" TEXT,
    "preferedarticlecategories" TEXT[],
    "profilepicture" TEXT
);
