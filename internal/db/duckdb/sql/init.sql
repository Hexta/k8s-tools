CREATE SCHEMA IF NOT EXISTS k8s;
CREATE TABLE IF NOT EXISTS k8s.nodes (
    name STRING,
    age BIGINT,
    creation_ts TIMESTAMP,
    instance_type STRING,
    cpu_utilisation FLOAT,
    memory_utilisation FLOAT
);
