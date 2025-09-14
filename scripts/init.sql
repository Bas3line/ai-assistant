-- Development database initialization script
-- This script is run when the PostgreSQL container starts for the first time

-- Create the database if it doesn't exist (it should already exist from POSTGRES_DB)
-- CREATE DATABASE IF NOT EXISTS ai_assistant;

-- Connect to the ai_assistant database
\c ai_assistant;

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create a test user for development (optional)
-- You can remove this in production
CREATE ROLE ai_dev_user WITH LOGIN PASSWORD 'dev_password';
GRANT ALL PRIVILEGES ON DATABASE ai_assistant TO ai_dev_user;

-- Log the initialization
\echo 'Database ai_assistant initialized successfully';