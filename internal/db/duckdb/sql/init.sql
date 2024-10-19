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

CREATE TABLE IF NOT EXISTS k8s.taints (
    node_name STRING,
    effect STRING,
    key STRING,
    value STRING,
    PRIMARY KEY (node_name, key, value, effect)
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

CREATE TABLE IF NOT EXISTS k8s.services (
    name STRING,
    namespace STRING,
    labels MAP(STRING, STRING),
    annotations MAP(STRING, STRING),
    creation_ts TIMESTAMP,

    cluster_ip STRING,
    cluster_ips STRING[],
    external_ips STRING[],
    external_traffic_policy STRING,
    health_check_node_port INTEGER,
    load_balancer_class STRING,
    load_balancer_source_ranges STRING[],
    publish_not_ready_addresses BOOLEAN,
    reason STRING,
    selector MAP(STRING, STRING),
    session_affinity STRING,
    type STRING,
    PRIMARY KEY (name, namespace)
)
