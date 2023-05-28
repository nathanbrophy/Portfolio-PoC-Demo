# Define the provider for the TF AWS modules
provider "aws" {
    region = var.region
}

# Dynamically load the availability zone data
# from the specified var.region value. 
#
# These calues are defined later on in the module
# loading on the dynamic object parsing.
#
# We define the slice with an end length of 1 because we
# only have data in a single az and subnet. 
data "aws_availability_zones" "available" {}

# https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/latest
# Full documentation on this module is linked above, the module is used
# toprovision the VPC and subnet definitions for the cluster environment.
module "vpc" {
    source  = "terraform-aws-modules/vpc/aws"
    version = "3.19.0"

    name = "example-vpc"

    cidr = var.vpc_cidr
    # This is set to a 2 len slice with 2 subnets for private/public
    # because it is a hard requirement from the TF provider module to
    # have more than one subnet AZ for the created VPC
    azs = slice(data.aws_availability_zones.available.names, 0, 2)

    private_subnets = var.private_subnets
    public_subnets  = var.public_subnets

    enable_nat_gateway   = true
    single_nat_gateway   = true
    enable_dns_hostnames = true

    public_subnet_tags = {
        "kubernetes.io/cluster/${var.cluster_name}" = "shared"
        "kubernetes.io/role/elb"                      = 1
    }

    private_subnet_tags = {
        "kubernetes.io/cluster/${var.cluster_name}" = "shared"
        "kubernetes.io/role/internal-elb"             = 1
    }
}

# https://registry.terraform.io/modules/terraform-aws-modules/eks/aws/latest
# The module is used to provision the EKS cluster 
# in the AWS account and environment. 
module "eks" {
    source  = "terraform-aws-modules/eks/aws"
    version = "19.5.1"

    cluster_name    = var.cluster_name
    cluster_version = var.cluster_version

    vpc_id                         = module.vpc.vpc_id
    subnet_ids                     = module.vpc.private_subnets
    cluster_endpoint_public_access = true

    eks_managed_node_group_defaults = {
        ami_type = "AL2_x86_64"
    }

    eks_managed_node_groups = {
        exmaple-node-group = {
            name = "example-node-group"

            instance_types = ["t3.small"]

            min_size     = 1
            max_size     = 1
            desired_size = 1
        }
    }
}
