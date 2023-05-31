output "oidc_issuer_id" {
    value = module.eks.oidc_provider
}

output "k8s_cluster" {
    value = module.eks.cluster_endpoint
}

output "k8s_ca_Data" {
    value = module.eks.cluster_certificate_authority_data
    sensitive = true
}

output "k8s_cluster_name" {
  value = module.eks.cluster_name
}
