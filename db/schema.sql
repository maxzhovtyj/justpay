CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS orders
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES users (id),
    status     INTEGER NOT NULL DEFAULT -1,
    is_final   bool             DEFAULT false,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS events
(
    id         uuid PRIMARY KEY,
    order_id   UUID REFERENCES orders (id),
    user_id    UUID REFERENCES users (id),
    status     INTEGER NOT NULL DEFAULT -1,
    is_final   bool             DEFAULT false,
    created_at timestamptz,
    updated_at timestamptz
);