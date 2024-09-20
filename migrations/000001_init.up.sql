-- Пример исправленного скрипта миграции 0001_initial.up.sql
CREATE TABLE IF NOT EXISTS data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mark TEXT NOT NULL,
    model TEXT NOT NULL,
    year INTEGER,
    volume INTEGER,
    amount INTEGER,
    popular_rate INTEGER
);

CREATE TABLE IF NOT EXISTS delivered (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    from_city TEXT,
    to_city TEXT,
    amount INTEGER,
    country TEXT
);

CREATE TABLE IF NOT EXISTS broker_amount (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    amount INTEGER,
    country TEXT
);

-- Вставка данных
INSERT INTO delivered (from_city, to_city, amount, country) VALUES ('Дубай', 'Алматы', 2500, 'kz');
INSERT INTO delivered (from_city, to_city, amount, country) VALUES ('Дубай', 'Шымкент', 2600, 'kz');
INSERT INTO delivered (from_city, to_city, amount, country) VALUES ('Дубай', 'Актау', 2000, 'kz');
INSERT INTO delivered (from_city, to_city, amount, country) VALUES ('Дубай', 'Астрахань', 2200, 'ru');
INSERT INTO delivered (from_city, to_city, amount, country) VALUES ('Дубай', 'Бишкек', 2600, 'kg');

INSERT INTO broker_amount (amount, country) VALUES (25000, 'kz');
INSERT INTO broker_amount (amount, country) VALUES (30000, 'kz');
INSERT INTO broker_amount (amount, country) VALUES (35000, 'kz');
