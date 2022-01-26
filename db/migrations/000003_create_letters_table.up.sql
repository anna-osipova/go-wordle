CREATE TABLE IF NOT EXISTS letters (
  id SERIAL PRIMARY KEY,
  attempt_id INTEGER,
  letter VARCHAR(1),
  color VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  constraint fk_attempts_letters
     foreign key (attempt_id) 
     REFERENCES attempts (id)
);