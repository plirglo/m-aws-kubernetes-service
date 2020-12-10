define M_METADATA_CONTENT
labels:
  version: $(M_VERSION)
  name: AWS Kubernetes Service
  short: $(M_MODULE_SHORT)
  kind: infrastructure
  provider: aws
endef

define M_CONFIG_CONTENT
kind: $(M_MODULE_SHORT)-config
$(M_MODULE_SHORT):
  name: $(M_NAME)
  vpc_id: $(M_VPC_ID)
  region: $(M_REGION)
  subnet_ids: $(M_SUBNET_IDS)
  private_route_table_id: $(M_PRIVATE_ROUTE_TABLE_ID)
  disk_size: $(M_DISK_SIZE)
  autoscaler_scale_down_utilization_threshold: $(M_AUTOSCALER_SCALE_DOWN_UTILIZATION_THRESHOLD)
  ami_type: $(M_AMI_TYPE)
  ec2_ssh_key: $(M_EC2_SSH_KEY)
  worker_groups.name: $(M_WORKER_GROUPS_NAME)
  worker_groups.instance_type: $(M_WORKER_GROUPS_INSTANCE_TYPE)
  worker_groups.asg_desired_capacity: $(M_WORKER_GROUPS_ADG_DESIRED_CAPACITY)
  worker_groups.asg_min_size: $(M_WORKER_GROUPS_ASG_MIN_SIZE)
  worker_groups.asg_max_size: $(M_WORKER_GROUPS_ASG_MAX_SIZE)
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
