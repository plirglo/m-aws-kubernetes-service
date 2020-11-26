module "control_plane" {
  source       = "./modules/control_plane"
  name         = var.name
  cluster_name = local.cluster_name
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
  cluster_name  = local.cluster_name
  subnet_ids    = local.subnet_ids
  worker_groups = var.worker_groups
  depends_on    = [module.control_plane]
  providers = {
    aws = aws
  }
}

module "autoscaler" {
  source                   = "./modules/autoscaler"
  name                     = var.name
  cluster_name             = local.cluster_name
  region                   = var.region
  openid_connect_arn       = module.control_plane.openid_connect_arn
  openid_connect_url       = module.control_plane.openid_connect_url
  autoscaler_version       = local.autoscaler_version
  autoscaler_chart_version = "7.3.4"
  depends_on               = [module.control_plane, module.nodes]

  # https://discuss.hashicorp.com/t/module-does-not-support-depends-on/11692/3
  providers = {
    helm       = helm
    kubernetes = kubernetes
  }
}
