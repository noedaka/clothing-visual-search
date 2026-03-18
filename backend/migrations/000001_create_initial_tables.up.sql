CREATE TABLE categories (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        VARCHAR(255) NOT NULL,
);

CREATE TABLE products (
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    price        DECIMAL(10,2) NOT NULL,
    category_id  INT REFERENCES categories(id) ON DELETE SET NULL
);

CREATE TABLE product_images (
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    product_id   INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url          TEXT NOT NULL,
    is_primary   BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order   INT NOT NULL DEFAULT 0
);
