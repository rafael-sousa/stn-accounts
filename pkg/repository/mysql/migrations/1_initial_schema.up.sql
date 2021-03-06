CREATE TABLE account(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cpf CHAR(11) NOT NULL UNIQUE,
    secret CHAR(64) NOT NULL,
    balance BIGINT DEFAULT '0',
    created_at DATETIME NOT NULL
);

CREATE TABLE transfer(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    account_origin_id INT NOT NULL REFERENCES account(id),
    account_destination_id INT NOT NULL REFERENCES account(id),
    amount BIGINT DEFAULT '0' CHECK(amount > 0),
    created_at DATETIME NOT NULL,
    CONSTRAINT origin_dest CHECK (account_origin_id != account_destination_id)
);