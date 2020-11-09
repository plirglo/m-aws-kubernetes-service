data "aws_vpc" "vpc" {
  id = var.vpc_id
}

data "aws_route_table" "private_route_table" {
  count          = var.subnet_ids != null ? 0 : 1
  route_table_id = var.private_route_table_id
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "eks_subnet" {
  # Subnets in at least 2 availability zones are required for EKS
  count             = var.subnet_ids != null ? 0 : 2
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(data.aws_vpc.vpc.cidr_block, 4, 15-count.index)
  vpc_id            = data.aws_vpc.vpc.id
  tags              = {
    name                                    = "${var.name}-eks-subnet${count.index}"
    cluster_name                            = var.name
    # https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html#vpc-subnet-tagging
    "kubernetes.io/cluster/${var.name}-eks" = "shared"
    "kubernetes.io/role/internal-elb"       = 1
  }
}

resource "aws_route_table_association" "private" {
  count          = var.subnet_ids != null ? 0 : 2
  subnet_id      = aws_subnet.eks_subnet[count.index].id
  route_table_id = data.aws_route_table.private_route_table[0].route_table_id
}

# https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html#vpc-tagging
resource "aws_ec2_tag" "eks_vpc" {
  count       = var.subnet_ids != null ? 0 : 1
  resource_id = data.aws_vpc.vpc.id
  key         = "kubernetes.io/cluster/${module.control_plane.cluster_name}"
  value       = "shared"
}
