CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    title TEXT,
    amount FLOAT,
    note TEXT,
    tags TEXT[]
);

INSERT INTO expenses (id, title, amount, note, tags) VALUES (DEFAULT, 'Bob', 20, 'testing', ARRAY['foo', 'bar']);