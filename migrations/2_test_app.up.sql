INSERT INTO apps("name", "secret")
VALUES ('test', 'test')
ON CONFLICT DO NOTHING;