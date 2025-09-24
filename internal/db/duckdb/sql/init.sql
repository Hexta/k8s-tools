CREATE SCHEMA IF NOT EXISTS k8s;

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

CREATE TABLE IF NOT EXISTS k8s.daemonsets (
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
CREATE VIEW IF NOT EXISTS k8s.ds AS SELECT * FROM k8s.daemonsets;

CREATE TABLE IF NOT EXISTS k8s.endpoint_slices (
    name STRING,
    namespace STRING,
    annotations MAP(STRING, STRING),
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    address_type STRING,
    api_version STRING,
    kind STRING,
    ports STRING[],
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.endpoints (
    endpoint_slice_name STRING,
    namespace STRING,
    addresses STRING[],
    hostname STRING,
    ip STRING,
    is_ready BOOLEAN,
    node_name STRING,
    has_target_ref BOOLEAN,
    target_ref_api_version STRING,
    target_ref_field_path STRING,
    target_ref_kind STRING,
    target_ref_name STRING,
    target_ref_namespace STRING,
    target_ref_resource_version STRING,
    target_ref_uid STRING
);

CREATE TABLE IF NOT EXISTS k8s.endpoint_slice_ports (
    endpoint_slice_name STRING,
    namespace STRING,
    app_protocol STRING,
    name STRING,
    port INTEGER,
    protocol STRING,
    PRIMARY KEY (namespace, endpoint_slice_name, name)
);

CREATE TABLE IF NOT EXISTS k8s.horizontal_pod_autoscalers (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    current_replicas INTEGER,
    desired_replicas INTEGER,
    PRIMARY KEY (namespace, name)
);
CREATE VIEW IF NOT EXISTS k8s.hpa AS SELECT * FROM k8s.horizontal_pod_autoscalers;

CREATE TABLE IF NOT EXISTS k8s.images (
    names STRING[],
    node_name STRING,
    size INTEGER
);

CREATE TABLE IF NOT EXISTS k8s.init_containers (
    name STRING,
    namespace STRING,
    pod_name STRING,
    cpu_requests FLOAT,
    cpu_limits FLOAT,
    memory_requests FLOAT,
    memory_limits FLOAT,
    PRIMARY KEY (namespace, pod_name, name)
);

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

CREATE TABLE IF NOT EXISTS k8s.persistent_volume_claims (
    name STRING,
    namespace STRING,
    labels MAP(STRING, STRING),
    annotations MAP(STRING, STRING),
    creation_ts TIMESTAMP,

    access_modes STRING[],
    allocated_resource_statuses MAP(STRING, STRING),
    allocated_resources MAP(STRING, INT64),
    capacity MAP(STRING, INT64),
    current_volume_attributes_class_name STRING,
    desired_access_modes STRING[],
    modify_volume_status STRING,
    phase STRING,
    resource_limits MAP(STRING, INT64),
    resource_requests MAP(STRING, INT64),
    storage_class_name STRING,
    target_volume_attributes_class_name STRING,
    volume_attributes_class_name STRING,
    volume_mode STRING,
    volume_name STRING,
    PRIMARY KEY (namespace, name)
);
CREATE VIEW IF NOT EXISTS k8s.pvcs AS SELECT * FROM k8s.persistent_volume_claims;

CREATE TABLE IF NOT EXISTS k8s.persistent_volumes (
    name STRING,
    labels MAP(STRING, STRING),
    annotations MAP(STRING, STRING),
    creation_ts TIMESTAMP,

    access_modes STRING[],
    capacity MAP(STRING, INT64),
    claim_ref_kind STRING,
    claim_ref_name STRING,
    claim_ref_namespace STRING,
    persistent_volume_reclaim_policy STRING,
    phase STRING,
    storage_class_name STRING,
    volume_mode STRING,
    PRIMARY KEY (name)
);
CREATE VIEW IF NOT EXISTS k8s.pvs AS SELECT * FROM k8s.persistent_volumes;

CREATE TABLE IF NOT EXISTS k8s.persistent_volume_claim_volumes (
    name STRING,
    pod_name STRING,
    namespace STRING,

    claim_name STRING,
    read_only BOOLEAN,

    PRIMARY KEY (namespace, pod_name, claim_name, name)
);
CREATE VIEW IF NOT EXISTS k8s.pvc_volumes AS SELECT * FROM k8s.persistent_volume_claim_volumes;

CREATE TABLE IF NOT EXISTS k8s.pods (
    name STRING,
    namespace STRING,
    labels MAP(STRING, STRING),
    annotations MAP(STRING, STRING),
    creation_ts TIMESTAMP,

    automount_service_account_token BOOLEAN,
    cpu_limits FLOAT,
    cpu_requests FLOAT,
    dns_policy STRING,
    enable_service_links BOOLEAN,
    host_ipc BOOLEAN,
    host_network BOOLEAN,
    host_pid BOOLEAN,
    host_users BOOLEAN,
    hostname STRING,
    ip STRING,
    memory_limits FLOAT,
    memory_requests FLOAT,
    node_name STRING,
    node_selector MAP(STRING, STRING),
    priority INTEGER,
    priority_class_name STRING,
    restart_policy STRING,
    runtime_class_name STRING,
    scheduler_name STRING,
    service_account_name STRING,
    set_hostname_as_fqdn BOOLEAN,
    share_process_namespace BOOLEAN,
    subdomain STRING,
    termination_grace_period_seconds INTEGER,
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
    PRIMARY KEY (namespace, name)
);

CREATE TABLE IF NOT EXISTS k8s.stateful_sets (
    name STRING,
    namespace STRING,
    creation_ts TIMESTAMP,
    labels MAP(STRING, STRING),
    replicas INTEGER,
    PRIMARY KEY (namespace, name)
);
CREATE VIEW IF NOT EXISTS k8s.sts AS SELECT * FROM k8s.stateful_sets;

CREATE TABLE IF NOT EXISTS k8s.taints (
    node_name STRING,
    effect STRING,
    key STRING,
    value STRING,
    PRIMARY KEY (node_name, key, value, effect)
);

CREATE TABLE IF NOT EXISTS k8s.tolerations (
    pod_name STRING,
    effect STRING,
    key STRING,
    operator STRING,
    toleration_seconds INTEGER,
    value STRING,
);
