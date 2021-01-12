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
      name                 = "default_wg_lin"
      instance_type        = "t2.small"
      asg_desired_capacity = 1
      asg_min_size         = 1
      asg_max_size         = 1
    }
  ]
}

variable "worker_groups_win" {
  description = "Worker groups definition list - Windows"
  type        = list(object({
    name                 = string
    instance_type        = string
    asg_desired_capacity = number
    asg_min_size         = number
    asg_max_size         = number
  }))
  default     = [
    {
      name                 = "default_wg_win"
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
