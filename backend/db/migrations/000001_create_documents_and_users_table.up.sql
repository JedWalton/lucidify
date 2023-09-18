-- Users table
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY, -- Changed from SERIAL to VARCHAR to store the provided ID
    external_id VARCHAR(255), -- New column
    username VARCHAR(255) UNIQUE, -- Kept as is, but might be nullable since the webhook has null values
    password_enabled BOOLEAN, -- New column to store if password is enabled
    email VARCHAR(255) UNIQUE NOT NULL, -- Kept as is
    first_name VARCHAR(255), -- New column
    last_name VARCHAR(255), -- New column
    image_url TEXT, -- New column
    profile_image_url TEXT, -- New column
    two_factor_enabled BOOLEAN, -- New column
    created_at TIMESTAMP, -- Adjusted to not have default value
    updated_at TIMESTAMP, -- Adjusted to not have default value
    deleted BOOLEAN DEFAULT FALSE -- New column to track if a user is deleted
);

-- Documents table
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, -- Changed data type to VARCHAR(255)
    document_name VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, document_name) -- Composite unique constraint remains unchanged
);

-- Index on user_id for the documents table remains unchanged
CREATE INDEX idx_documents_user_id ON documents(user_id);
