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
    cpu_limits FLOAT,
    memory_requests FLOAT,
    memory_limits FLOAT,
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.containers (
    name STRING,
    namespace STRING,
    pod STRING,
    cpu_requests FLOAT,
    cpu_limits FLOAT,
    memory_requests FLOAT,
    memory_limits FLOAT,
    PRIMARY KEY (namespace, pod, name)
);

CREATE TABLE IF NOT EXISTS k8s.deployments (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    PRIMARY KEY (namespace, name)
);
