tls_mode                   = "STRICT"
tls_enforce                = true
port                       = 8080
health_check_path          = "/"
certificate_authority_arns = []
logical_product_family     = "terratest"
logical_product_service    = "vservicetest"
tags = {
  "env" : "gotest",
  "creator" : "terratest",
  "provisioner" : "Terraform",
}
