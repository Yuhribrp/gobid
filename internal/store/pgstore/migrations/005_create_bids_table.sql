CREATE TABLE IF NOT EXISTS bids (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    bidder_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    create_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    bid_amount FLOAT8 NOT NULL
);

---- create above / drop below ----

DROPCON TABLE IF EXISTS bids;
