CREATE TABLE addresses (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `host_id` INTEGER NOT NULL,
    `addr` TEXT NOT NULL,
    `addrtype` TEXT,
    FOREIGN KEY(`host_id`) REFERENCES hosts(id)
);