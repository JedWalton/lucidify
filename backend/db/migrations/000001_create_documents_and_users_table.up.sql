-- Users table
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    external_id VARCHAR(255),
    username VARCHAR(255) UNIQUE,
    password_enabled BOOLEAN,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    image_url TEXT,
    profile_image_url TEXT,
    two_factor_enabled BOOLEAN,
    created_at BIGINT, -- Changed to BIGINT
    updated_at BIGINT, -- Changed to BIGINT
    deleted BOOLEAN DEFAULT FALSE
);

-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE documents (
    document_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    document_name VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, document_name),
    UNIQUE(user_id, document_id)  -- Add this line
);

-- Index on user_id for the documents table remains unchanged
CREATE INDEX idx_documents_user_id ON documents(user_id);

-- Chunks table with user_id FK
CREATE TABLE document_chunks (
    chunk_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    document_id UUID NOT NULL REFERENCES documents(document_id) ON DELETE CASCADE,
    chunk_content TEXT NOT NULL,
    chunk_index INT NOT NULL, -- To keep track of the order of chunks for a document
    FOREIGN KEY (user_id, document_id) REFERENCES documents(user_id, document_id) ON DELETE CASCADE
);

-- Index on document_id for the document_chunks table
CREATE INDEX idx_document_chunks_document_id ON document_chunks(document_id);
-- Index on user_id for the document_chunks table
CREATE INDEX idx_document_chunks_user_id ON document_chunks(user_id);

