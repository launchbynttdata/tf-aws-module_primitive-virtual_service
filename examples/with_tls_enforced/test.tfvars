naming_prefix              = "demo"
tls_mode                   = "STRICT"
tls_enforce                = true
port                       = 8080
health_check_path          = "/"
certificate_authority_arns = []
tags = {
  "env" : "gotest",
  "creator" : "terratest",
  "provisioner" : "Terraform",
}
