CREATE TABLE IF NOT EXISTS attempts (
  id SERIAL PRIMARY KEY,
  session_id uuid,
  word VARCHAR(5),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  constraint fk_sessions_attempts
     foreign key (session_id) 
     REFERENCES sessions (id)
);