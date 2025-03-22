module "network" {
  source = "./modules/network"
}

module "secret" {
  source = "./modules/secret"
}

module "login" {
  source = "./modules/login"

  lambda_name = var.lambda_name
  vpc_name    = var.vpc_name

  sign_key = module.secret.sign_key

  security_group_id = module.network.security_group_id

  depends_on = [
    module.secret
  ]
}