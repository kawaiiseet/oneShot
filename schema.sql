CREATE TABLE IF NOT EXISTS personal_data
(
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    login         TEXT        NOT NULL UNIQUE,
    age           INT         NOT NULL,
    country       TEXT        NOT NULL,
    interests     TEXT[]      NOT NULL DEFAULT '{}',
    likes         TEXT[]      NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_personal_data_country ON personal_data (country);
CREATE INDEX IF NOT EXISTS idx_personal_data_interests ON personal_data USING GIN (interests);
