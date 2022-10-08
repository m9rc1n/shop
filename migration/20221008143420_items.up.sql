CREATE SCHEMA IF NOT EXISTS shop;

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS shop.items
(
    id             SERIAL       NOT NULL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL UNIQUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    quantity       INTEGER      NOT NULL,
    reserved_at    TIMESTAMPTZ,
    reservation_id BIGINT
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON shop.items
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

