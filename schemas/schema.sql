CREATE TABLE customer
(
    customerId SERIAL,
    customerEmail TEXT NOT NULL,
    customerPhone VARCHAR(12) NOT NULL,
    customerParcelWeight NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT customer_pkey PRIMARY KEY (customerId)
)
