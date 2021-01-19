module "control_plane" {
  source       = "./modules/control_plane"
  name         = var.name
  k8s_version  = var.k8s_version
  subnet_ids   = local.subnet_ids
  providers    = {
    aws      = aws
    tls      = tls
    template = template
  }
}

module "nodes" {
  source        = "./modules/nodes"
  name          = var.name
  subnet_ids    = local.subnet_ids
  worker_groups = var.worker_groups
  depends_on    = [module.control_plane]
  providers     = {
    aws = aws
  }
}

module "autoscaler" {
  source                   = "./modules/autoscaler"
  name                     = var.name
  region                   = var.region
  openid_connect_arn       = module.control_plane.openid_connect_arn
  openid_connect_url       = module.control_plane.openid_connect_url
  autoscaler_version       = local.autoscaler_version
  autoscaler_chart_version = "7.3.4"
  depends_on               = [module.control_plane, module.nodes]

  # https://discuss.hashicorp.com/t/module-does-not-support-depends-on/11692/3
  providers                = {
    helm       = helm
    kubernetes = kubernetes
  }
}

resource "null_resource" "eks_nodes_win" {
   provisioner "local-exec" {
      command =<<EOF
      "eksctl create nodegroup \
         --cluster var.name \
         --region var.worker_groups_win.region \
         --name var.worker_groups_win.name \
         --node-type var.worker_groups_win.instance_type \
         --nodes var.worker_groups_win.asg_desired_capacity \
         --nodes-min var.worker_groups_win.asg_min_size \
         --nodes-max var.worker_groups_win.asg_max_size \
         --node-ami-family var.worker_groups_win.platform"
         EOF

   }
}