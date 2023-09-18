-- Create the documents table
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    document_name VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
