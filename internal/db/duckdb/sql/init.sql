CREATE SCHEMA IF NOT EXISTS k8s;
CREATE TABLE IF NOT EXISTS k8s.nodes (
    name STRING PRIMARY KEY,
    creation_ts TIMESTAMP,
    instance_type STRING,
    cpu_utilisation FLOAT,
    memory_utilisation FLOAT,
    labels MAP(STRING, STRING)
);

CREATE TABLE IF NOT EXISTS k8s.pods (
    name STRING,
    namespace STRING,
    node STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    cpu_requests FLOAT,
    memory_requests FLOAT,
    PRIMARY KEY (namespace, name)
);
