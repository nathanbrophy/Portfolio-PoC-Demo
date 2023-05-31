variable "region" {
    description = "AWS region"
    type        = string
    default     = "us-west-1"
}

variable "cluster_name" {
    description = "Cluster Name"
    type        = string
    default     = "helmApplyTF"
}

variable "cluster_version" {
    description = "Cluster Name"
    type        = string
    default     = "1.25"
}

variable "vpc_cidr" {
    description = "VPC CIDR"
    type        = string
    default     = "10.0.0.0/16"
}

variable "private_subnets" {
    description = "Private Subnet CIDR Ranges"
    type        = list(string)
    default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "public_subnets" {
    description = "Public Subnet CIDR Ranges"
    type        = list(string)
    default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "iam_policy_name" {
    description = "IAM policy name"
    type        = string
    default     = "ingress_policy"
}

variable "iam_role_name" {
    description = "IAM role name"
    type        = string
    default     = "ingress_alb_controller"
}
