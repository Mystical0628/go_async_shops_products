CREATE TABLE shops (
    id INT(11) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    url VARCHAR(255) NOT NULL,
    opens_at TIME NOT NULL,
    closes_at TIME NOT NULL
);

ALTER TABLE products ADD
    CONSTRAINT products_shops_fk
    FOREIGN KEY (shop_id)
        REFERENCES shops (id);