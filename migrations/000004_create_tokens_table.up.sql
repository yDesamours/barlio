CREATE TABLE IF NOT EXISTS tokens(
    "userid" INT references "users"("id") not null,
    "hash" BYTEA NOT NULL,
    "scope" TEXT NOT NULL,
    "expire_at" TIMESTAMP NOT NULL
);

ALTER TABLE "tokens" ADD PRIMARY KEY("userid", "scope");