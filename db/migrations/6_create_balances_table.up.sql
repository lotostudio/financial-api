CREATE TABLE IF NOT EXISTS balances(
    account_id BIGINT NOT NULL,
    date DATE NOT NULL,
    value NUMERIC NOT NULL,
    CONSTRAINT fk_balance_account FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    PRIMARY KEY(account_id, date)
);

-- create new entries for every account
INSERT INTO balances(account_id, date, value)
SELECT id, current_date, balance FROM accounts
