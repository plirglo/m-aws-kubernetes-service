
resource "null_resource" "install_vpc_controller" {
   provisioner "local-exec" {
      command = "eksctl utils install-vpc-controllers --cluster ${name} --approve" 
}

module "eks" {
  source          = "../.."
  cluster_name    = var.name
  subnets         = var.subnet_ids
  vpc_id          = var.vpc_id

  tags = {
      tags = local.eks_node_tags
  }

  worker_groups = [
    {
      name                          = "worker-group-windows"
      instance_type                 = "t2.small"
      platform                      = "windows"
      asg_desired_capacity          = 1
    }
  ]
}