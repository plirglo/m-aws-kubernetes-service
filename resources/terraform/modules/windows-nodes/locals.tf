locals {
  tags = map(
    "resource_group", var.name
  )
  eks_node_tags = merge(
    local.tags,
    map(
      "k8s.io/cluster-autoscaler/enabled", "true",
      "k8s.io/cluster-autoscaler/${var.name}", "true"
    )
  )
}
