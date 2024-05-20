BEGIN;

CREATE TYPE "OrderStatus" AS ENUM ('PENDING', 'COMPLETED', 'CANCELLED');

CREATE TABLE IF NOT EXISTS orders (
  id SERIAL NOT NULL PRIMARY KEY,
  "userId" INTEGER NOT NULL,
  total DOUBLE PRECISION NOT NULL,
  status "OrderStatus" NOT NULL DEFAULT 'PENDING',
  address TEXT NOT NULL,
  "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_user FOREIGN KEY("userId") REFERENCES users(id)
);

COMMIT;