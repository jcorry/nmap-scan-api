CREATE TABLE hostnames(
    `id` INTEGER NOT NULL PRIMARY KEY,
    `host_id` INTEGER NOT NULL,
    `name` TEXT NOT NULL,
    `type` TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);