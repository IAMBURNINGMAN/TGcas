-- тестовый промокод на 1000 Сабинок
-- запустить один раз после первого старта бота:
-- psql $DATABASE_URL -f migrations/seed_test_promo.sql

INSERT INTO promo_codes (code, amount) VALUES ('SABINKA1000', 1000)
ON CONFLICT (code) DO NOTHING;
