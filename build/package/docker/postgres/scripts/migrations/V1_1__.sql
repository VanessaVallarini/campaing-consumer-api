CREATE TABLE IF NOT EXISTS campaing (
  id uuid PRIMARY KEY NOT NULL,
  user_id uuid NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp,
  slug VARCHAR(10) NOT NULL,
  active BOOLEAN DEFAULT true,
  lat DECIMAL NOT NULL,
  long DECIMAL NOT NULL,
  qtd_max INTEGER,
  qtd_min INTEGER,
  position INTEGER,
);
