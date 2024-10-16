CREATE SCHEMA IF NOT EXISTS k8s;

CREATE TABLE IF NOT EXISTS k8s.nodes (
    name STRING PRIMARY KEY,
    labels MAP(STRING, STRING),

    address MAP(STRING, STRING),
    allocatable_cpu FLOAT,
    allocatable_memory FLOAT,
    annotations MAP(STRING, STRING),
    architecture STRING,
    capacity_cpu FLOAT,
    capacity_memory FLOAT,
    container_runtime_version STRING,
    cpu_utilisation FLOAT,
    creation_ts TIMESTAMP,
    instance_type STRING,
    kernel_version STRING,
    kubelet_version STRING,
    memory_utilisation FLOAT,
    operating_system STRING,
    os_image STRING
);

CREATE TABLE IF NOT EXISTS k8s.pods (
    name STRING,
    namespace STRING,
    node_name STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    cpu_requests FLOAT,
    cpu_limits FLOAT,
    memory_requests FLOAT,
    memory_limits FLOAT,
    ip STRING,
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.containers (
    name STRING,
    namespace STRING,
    pod_name STRING,
    cpu_requests FLOAT,
    cpu_limits FLOAT,
    memory_requests FLOAT,
    memory_limits FLOAT,
    PRIMARY KEY (namespace, pod_name, name)
);

CREATE TABLE IF NOT EXISTS k8s.deployments (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    replicas INTEGER,
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.hpa (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    current_replicas INTEGER,
    desired_replicas INTEGER,
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.sts (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    replicas INTEGER,
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.ds (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    current_number_scheduled INTEGER,
	desired_number_scheduled INTEGER,
	number_available INTEGER,
	number_misscheduled INTEGER,
	number_ready INTEGER,
	number_unavailable INTEGER,
    PRIMARY KEY (namespace, name)
);
