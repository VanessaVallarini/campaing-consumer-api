CREATE SCHEMA campaing;

GRANT ALL PRIVILEGES ON DATABASE "campaing-consumer-api-db" TO "postgres";

GRANT USAGE ON SCHEMA campaing TO "postgres";
ALTER USER "postgres" SET search_path = 'campaing';


SET SCHEMA 'campaing';
ALTER DEFAULT PRIVILEGES
    IN SCHEMA campaing
GRANT SELECT, UPDATE, INSERT, DELETE ON TABLES
    TO "postgres";

ALTER DEFAULT PRIVILEGES
    IN SCHEMA campaing
GRANT USAGE ON SEQUENCES
    TO "postgres";