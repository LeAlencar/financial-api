-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS financial
WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};

-- Use the keyspace
USE financial;

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    user_id UUID,
    transaction_type TEXT,
    amount DECIMAL,
    currency TEXT,
    status TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create validation_logs table with complete structure
CREATE TABLE IF NOT EXISTS validation_logs (
    id UUID PRIMARY KEY,
    event_type TEXT,
    action TEXT,
    event_id TEXT,
    user_id TEXT,
    status TEXT,
    message TEXT,
    details MAP<TEXT, TEXT>,
    raw_payload TEXT,
    processed_at TIMESTAMP,
    source TEXT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS ON transactions (user_id);
CREATE INDEX IF NOT EXISTS ON transactions (status);
CREATE INDEX IF NOT EXISTS ON validation_logs (event_type);
CREATE INDEX IF NOT EXISTS ON validation_logs (status);
CREATE INDEX IF NOT EXISTS ON validation_logs (user_id); 