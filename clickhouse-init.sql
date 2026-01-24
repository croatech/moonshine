-- Enable experimental MaterializedPostgreSQL
SET allow_experimental_database_materialized_postgresql = 1;

-- Create MaterializedPostgreSQL database for analytics
-- ClickHouse will automatically create replication slot
CREATE DATABASE IF NOT EXISTS moonshine_analytics
ENGINE = MaterializedPostgreSQL('postgres:5432', 'moonshine', 'postgres', 'postgres')
SETTINGS 
    materialized_postgresql_tables_list = 'movement_logs,rounds';

SELECT 'ClickHouse MaterializedPostgreSQL configured' as status;
