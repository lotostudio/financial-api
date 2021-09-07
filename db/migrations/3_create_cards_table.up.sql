CREATE TABLE IF NOT EXISTS cards(
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(4) NOT NULL,
    account_id BIGINT NOT NULL,
    CONSTRAINT oto_card_account FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- Create data entries in fresh table for all 'card' accounts
INSERT INTO cards(number, account_id)
SELECT '0000', a.id FROM accounts a WHERE a.type = 'card';