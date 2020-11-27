locals {
  k8s_service_account_namespace               = "kube-system"
  k8s_service_account_name                    = "cluster-autoscaler-aws-cluster-autoscaler"

  # https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-scale-down-work
  autoscaler_scale_down_utilization_threshold = "0.65"

  tags = map(
    "resource_group", var.name
  )
}
