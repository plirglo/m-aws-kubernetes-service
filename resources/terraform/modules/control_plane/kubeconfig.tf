data "template_file" "kubeconfig" {
  template = file("${path.module}/templates/kubeconfig.tpl")
  vars     = {
    endpoint     = aws_eks_cluster.eks_cluster.endpoint
    certificate  = aws_eks_cluster.eks_cluster.certificate_authority[0].data
    cluster_name = aws_eks_cluster.eks_cluster.name
  }
}
