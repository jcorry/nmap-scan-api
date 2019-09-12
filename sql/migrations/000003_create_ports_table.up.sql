CREATE TABLE ports (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `host_id` INTEGER NOT NULL,
    `protocol` TEXT NULL,
    `port_id` INTEGER NOT NULL,
    `owner` TEXT NULL,
    `service` TEXT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);