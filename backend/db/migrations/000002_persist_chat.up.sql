-- conversationHistory table
CREATE TABLE conversation_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) REFERENCES users(user_id),
    data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- folders table
CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) REFERENCES users(user_id),
    data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- prompts table
CREATE TABLE prompts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) REFERENCES users(user_id),
    data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger function to update 'updated_at'
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for 'conversation_history' table
CREATE TRIGGER tr_conversation_history_update_timestamp
BEFORE UPDATE ON conversation_history
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Trigger for 'folders' table
CREATE TRIGGER tr_folders_update_timestamp
BEFORE UPDATE ON folders
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Trigger for 'prompts' table
CREATE TRIGGER tr_prompts_update_timestamp
BEFORE UPDATE ON prompts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
