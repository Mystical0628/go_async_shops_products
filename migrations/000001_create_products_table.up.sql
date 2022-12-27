CREATE TABLE products (
    id INT(11) UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    shop_id INT(11) UNSIGNED NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(15,2) NOT NULL
);