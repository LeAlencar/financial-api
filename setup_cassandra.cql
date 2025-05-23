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

-- Create validation_logs table
CREATE TABLE IF NOT EXISTS validation_logs (
    id UUID PRIMARY KEY,
    transaction_id UUID,
    validation_type TEXT,
    result TEXT,
    message TEXT,
    created_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS ON transactions (user_id);
CREATE INDEX IF NOT EXISTS ON transactions (status);
CREATE INDEX IF NOT EXISTS ON validation_logs (transaction_id); 