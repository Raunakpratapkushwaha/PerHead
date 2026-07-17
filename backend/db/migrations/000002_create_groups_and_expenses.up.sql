CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members (
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE expenses (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT REFERENCES groups(id) ON DELETE CASCADE,
    payer_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL, -- Stored in cents (e.g., $10.00 = 1000)
    description VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'general',
    split_type VARCHAR(20) NOT NULL, -- 'EQUAL', 'EXACT', 'PERCENTAGE', 'SHARES'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE expense_splits (
    id BIGSERIAL PRIMARY KEY,
    expense_id BIGINT REFERENCES expenses(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL, -- What this specific user owes (in cents)
    percentage NUMERIC(5, 2) DEFAULT 0.00, -- Optional metadata (e.g., 33.33%)
    share INT DEFAULT 0, -- Optional metadata for share-based splits
    UNIQUE (expense_id, user_id)
);

-- Indexing for rapid balance lookups and ledger aggregates
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_expenses_group_id ON expenses(group_id);
CREATE INDEX idx_expense_splits_user_id ON expense_splits(user_id);
CREATE INDEX idx_expense_splits_expense_id ON expense_splits(expense_id);