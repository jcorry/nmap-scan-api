CREATE TABLE imports (
     id INTEGER NOT NULL PRIMARY KEY,
     file_id string,
     created string
);

CREATE UNIQUE INDEX idx_imports_file_id ON imports(file_id);

CREATE TABLE hosts (
    `id` INTEGER PRIMARY KEY,
    `file_id` TEXT NOT NULL,
    `starttime` TEXT NOT NULL,
    `endtime` TEXT NOT NULL,
    `comment` TEXT,
    `status` TEXT,
    FOREIGN KEY (`file_id`) REFERENCES imports(`file_id`)
);

CREATE TABLE addresses (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `host_id` INTEGER NOT NULL,
    `addr` TEXT NOT NULL,
    `addrtype` TEXT,
    FOREIGN KEY(`host_id`) REFERENCES hosts(`id`)
);

CREATE TABLE ports (
    `id` INTEGER NOT NULL PRIMARY KEY,
    `host_id` INTEGER NOT NULL,
    `protocol` TEXT NULL,
    `port_id` INTEGER NOT NULL,
    `owner` TEXT NULL,
    `service` TEXT NULL,
    FOREIGN KEY(`host_id`) REFERENCES hosts(`id`)
);

CREATE TABLE hostnames(
    `id` INTEGER NOT NULL PRIMARY KEY,
    `host_id` INTEGER NOT NULL,
    `name` TEXT NOT NULL,
    `type` TEXT NOT NULL,
    FOREIGN KEY(`host_id`) REFERENCES hosts(`id`)
);