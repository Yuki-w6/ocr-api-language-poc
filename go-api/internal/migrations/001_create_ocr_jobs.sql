CREATE TABLE IF NOT EXISTS ocr_jobs (
    id TEXT PRIMARY KEY,
    object_key TEXT NOT NULL,
    status TEXT NOT NULL,
    result_json JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ocr_jobs_status ON ocr_jobs (status);
CREATE INDEX IF NOT EXISTS idx_ocr_jobs_created_at ON ocr_jobs (created_at DESC);
