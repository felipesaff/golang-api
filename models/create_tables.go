package models

const (
	CREATE_TABLES = `
		CREATE TABLE IF NOT EXISTS users (
			id         BIGSERIAL PRIMARY KEY,
			name       TEXT NOT NULL,
			email      TEXT NOT NULL UNIQUE,
			document   TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS products (
			id          BIGSERIAL PRIMARY KEY,
			name        TEXT NOT NULL,
			description TEXT,
			price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
			stock       INT NOT NULL CHECK (stock >= 0),
			created_at  TIMESTAMP NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS orders (
			id         BIGSERIAL PRIMARY KEY,
			user_id    BIGINT NOT NULL REFERENCES users(id),
			total_cents BIGINT NOT NULL CHECK (total_cents >= 0),
			status     TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS order_items (
			order_id         BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id       BIGINT NOT NULL REFERENCES products(id),
			quantity         INT NOT NULL CHECK (quantity > 0),
			unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
			PRIMARY KEY (order_id, product_id)
		);


		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
	`
)
