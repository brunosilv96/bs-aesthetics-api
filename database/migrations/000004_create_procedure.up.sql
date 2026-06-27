CREATE TABLE IF NOT EXISTS procedures (
    id UUID PRIMARY KEY DEFAULT generate_uuid_v7(),
    registred_by UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    banner_url TEXT,
    price NUMERIC(10, 2) NOT NULL,
    duration_minutes INT NOT NULL,
    available BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);