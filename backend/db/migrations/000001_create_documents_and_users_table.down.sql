-- Drop indexes first
DROP INDEX IF EXISTS idx_document_chunks_user_id;
DROP INDEX IF EXISTS idx_document_chunks_document_id;
DROP INDEX IF EXISTS idx_documents_user_id;

-- Then drop the tables. Foreign key dependencies mean we drop the 'child' tables before the 'parent' tables
DROP TABLE IF EXISTS document_chunks;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS users;

-- If you are sure that no other tables in your database rely on the "uuid-ossp" extension, you can also remove the extension
DROP EXTENSION IF EXISTS "uuid-ossp";
