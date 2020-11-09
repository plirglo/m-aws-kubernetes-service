locals {
  tags = map(
    "k8s.io/cluster-autoscaler/enabled", "true",
    "k8s.io/cluster-autoscaler/${var.cluster_name}", "true"
  )
}
