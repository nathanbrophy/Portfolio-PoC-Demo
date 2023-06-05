resource "aws_iam_policy" "policy" {
    name = var.iam_policy_name
    description = var.description
    policy = var.policy
}

resource "aws_iam_role" "ingress_alb_controller" {
    name = var.iam_role_name

    force_detach_policies = true

    assume_role_policy = var.assume_role_policy
    
    tags = var.tags
}

resource "aws_iam_role_policy_attachment" "load-balancer-attach" {
  role       = aws_iam_role.ingress_alb_controller.name
  policy_arn = aws_iam_policy.policy.arn
}