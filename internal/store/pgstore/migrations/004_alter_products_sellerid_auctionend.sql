ALTER TABLE products
    ALTER COLUMN seller_id SET NOT NULL,
    ALTER COLUMN seller_id TYPE uuid USING seller_id::uuid,
    ALTER COLUMN auction_end TYPE timestamptz USING auction_end::timestamptz;

---- create above / drop below ----

ALTER TABLE products
    ALTER COLUMN seller_id DROP NOT NULL,
    ALTER COLUMN seller_id TYPE uuid USING seller_id::uuid,
    ALTER COLUMN auction_end TYPE timestamptz USING auction_end::timestamptz;
