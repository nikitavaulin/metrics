BEGIN;

CREATE TYPE metric_type AS ENUM (
    'gauge',
    'counter'
);

CREATE TABLE metrics (
    name TEXT PRIMARY KEY,
    mtype metric_type NOT NULL DEFAULT 'gauge',
    value DOUBLE PRECISION CHECK (mtype <> 'gauge' OR value IS NOT NULL),
    delta BIGINT CHECK (mtype <> 'counter' OR delta IS NOT NULL),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW; 
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER metrics_set_updated_at
BEFORE UPDATE ON metrics
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMIT;