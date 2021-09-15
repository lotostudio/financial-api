CREATE TYPE transaction_type AS ENUM('income', 'expense', 'transfer');

CREATE TABLE IF NOT EXISTS transaction_categories(
    id SERIAL PRIMARY KEY ,
    title VARCHAR(49) NOT NULL,
    type transaction_type NOT NULL
);


CREATE TABLE IF NOT EXISTS transactions(
    id BIGSERIAL PRIMARY KEY,
    amount NUMERIC NOT NULL,
    type transaction_type NOT NULL,
    created_at DATE NOT NULL,
    category_id INT,
    credit_id BIGINT,
    debit_id BIGINT,
    CONSTRAINT fk_transaction_category FOREIGN KEY(category_id) REFERENCES transaction_categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_transaction_credit FOREIGN KEY(credit_id) REFERENCES accounts(id) ON DELETE SET NULL,
    CONSTRAINT fk_transaction_debit FOREIGN KEY(debit_id) REFERENCES accounts(id) ON DELETE SET NULL
);