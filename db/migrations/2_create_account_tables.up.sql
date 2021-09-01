CREATE TABLE IF NOT EXISTS currencies(
    id SERIAL PRIMARY KEY ,
    code VARCHAR(10) NOT NULL
);

CREATE TYPE account_type AS ENUM('cash', 'card', 'loan', 'deposit');

CREATE TABLE IF NOT EXISTS accounts(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    balance NUMERIC NOT NULL,
    type account_type NOT NULL,
    currency_id INT NOT NULL,
    owner_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_account_owner FOREIGN KEY(currency_id) REFERENCES currencies(id) ON DELETE SET NULL,
    CONSTRAINT fk_account_currency FOREIGN KEY(owner_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS loans(
    id BIGSERIAL PRIMARY KEY,
    term INT NOT NULL,
    rate NUMERIC NOT NULL,
    account_id BIGINT NOT NULL,
    CONSTRAINT oto_loan_account FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS deposits(
    id BIGSERIAL PRIMARY KEY,
    term INT NOT NULL,
    rate NUMERIC NOT NULL,
    account_id BIGINT NOT NULL,
    CONSTRAINT oto_deposit_account FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);