M_NAME ?= epiphany
M_VPC_ID ?= unset
M_SUBNET_IDS ?= null
M_REGION ?= eu-central-1
M_PRIVATE_ROUTE_TABLE_ID ?= unset
M_DISK_SIZE ?= 36
M_AUTOSCALER_SCALE_DOWN_UTILIZATION_THRESHOLD ?= 0.65
M_EC2_SSH_KEY ?= null
M_AMI_TYPE ?= AL2_x86_64

define _M_WORKER_GROUPS
{
  name: default_wg,
  instance_type: t2.small,
  asg_desired_capacity: 1,
  asg_min_size: 1,
  asg_max_size: 1,
}
endef

M_WORKER_GROUPS ?= $(_M_WORKER_GROUPS)

# aws credentials
M_AWS_ACCESS_KEY ?= unset
M_AWS_SECRET_KEY ?= unset
