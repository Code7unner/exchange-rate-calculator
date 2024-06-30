CREATE TABLE exchange_rates (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    from_currency VARCHAR(10),
    to_currency VARCHAR(10),
    rate DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE exchange_rates ADD CONSTRAINT unique_currency_pair UNIQUE (from_currency, to_currency);

CREATE INDEX idx_exchange_rates_from_to ON exchange_rates (from_currency, to_currency);
