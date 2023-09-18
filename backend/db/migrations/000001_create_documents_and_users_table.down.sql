-- Drop the index first
DROP INDEX IF EXISTS idx_documents_user_id;

-- Drop the documents table next (because it has a foreign key reference to users)
DROP TABLE IF EXISTS documents;

-- Finally, drop the users table
DROP TABLE IF EXISTS users;

