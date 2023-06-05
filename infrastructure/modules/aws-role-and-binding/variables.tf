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

variable "policy" {
    description = "This is the IAM policy data to set on the created policy"
    type = string
}

variable "description" {
    description = "Description for the IAM policy object"
    type = string
}

variable "assume_role_policy" {
    description = "This is the JSON assume role data for the policy attachment to work properly"
    type = string
}

variable "tags" {
    description = "Additional tags to set on the created role"
    type = map(string)
    default = {}
}
