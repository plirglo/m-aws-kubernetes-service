variable "subnet_ids" {
  description = "Existing subnet ids to join to"
  type = list(string)
}

variable "name" {
  description = "Prefix for resource names"
  type        = string
  default     = "default"
}

variable "k8s_version" {
  description = "Kubernetes version to install"
  type        = string
  default     = "1.18"
}

variable "autoscaler_version" {
  description = "Kubernetes autoscaler version"
  type        = string
  default     = null
}

variable "vpc_id" {
  description = "VPC id to join to"
  type        = string
}

variable "worker_groups" {
  description = "Worker groups definition list"
  type        = list(object({
    name                 = string
    instance_type        = string
    asg_desired_capacity = number
    asg_min_size         = number
    asg_max_size         = number
  }))
  default     = [
    {
      name                 = "default_wg"
      instance_type        = "t2.small"
      asg_desired_capacity = 1
      asg_min_size         = 1
      asg_max_size         = 1
    }
  ]
}

variable "region" {
  description = "Region for AWS resources"
  type        = string
}

variable "private_route_table_id" {
  description = "Private route table id for table associations"
  type        = string
}

variable "disk_size" {
  description = "Disk size"
  type        = number
}

variable "autoscaler_scale_down_utilization_threshold" {
  description = "Autoscaler scale down utilization threshold"
  type        = string
}

variable "ami_type" {
  description = "Type of Amazon Machine Image (AMI) associated with the EKS Node Group"
  type        = string
}

variable "ec2_ssh_key" {
  description = "EC2 Key Pair name that provides access for SSH communication with the worker nodes in the EKS Node Group"
  type        = string
}
