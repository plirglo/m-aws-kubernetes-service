provider "aws" {
  region = var.region
}

provider "kubernetes" {
  host                   = module.control_plane.cluster_endpoint
  load_config_file       = false
  token                  = module.control_plane.cluster_token
  cluster_ca_certificate = base64decode(module.control_plane.cluster_ca)
}

provider "helm" {
  kubernetes {
    host                   = module.control_plane.cluster_endpoint
    load_config_file       = false
    token                  = module.control_plane.cluster_token
    cluster_ca_certificate = base64decode(module.control_plane.cluster_ca)
  }
}

provider "template" {}

provider "tls" {}
