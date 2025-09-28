terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.30"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

locals {
  name_prefix = "ssk-${var.environment}"
  tags = merge({
    System  = "ssk"
    Service = "uptime-monitor"
    Env     = var.environment
    Owner   = "ssk"
  }, var.additional_tags)
}

module "networking" {
  source      = "../../modules/networking"
  name_prefix = local.name_prefix
  cidr_block  = var.vpc_cidr
  azs         = var.availability_zones
  tags        = local.tags
}

module "dynamodb" {
  source      = "../../modules/dynamodb"
  name_prefix = local.name_prefix
  tags        = local.tags
}

module "iam" {
  source                   = "../../modules/iam"
  name_prefix              = local.name_prefix
  dynamodb_table_arns      = [module.dynamodb.sites_table_arn, module.dynamodb.status_table_arn]
  github_oidc_provider_arn = var.github_oidc_provider_arn
  github_repo              = var.github_repo
  tags                     = local.tags
}

module "ecs" {
  source                  = "../../modules/ecs_service"
  name_prefix             = local.name_prefix
  container_image         = var.container_image
  desired_count           = var.desired_count
  cpu                     = var.task_cpu
  memory                  = var.task_memory
  private_subnet_ids      = module.networking.private_subnet_ids
  public_subnet_ids       = module.networking.public_subnet_ids
  alb_security_group_id   = module.networking.alb_security_group_id
  ecs_security_group_id   = module.networking.ecs_security_group_id
  vpc_id                  = module.networking.vpc_id
  certificate_arn         = var.certificate_arn
  test_listener_port      = var.test_listener_port
  task_execution_role_arn = module.iam.task_execution_role_arn
  task_role_arn           = module.iam.task_role_arn
  tags                    = local.tags
}

module "observability" {
  source              = "../../modules/observability"
  name_prefix         = local.name_prefix
  cluster_name        = module.ecs.cluster_name
  service_name        = module.ecs.service_name
  target_group_arn    = module.ecs.blue_target_group_arn
  load_balancer_arn   = module.ecs.load_balancer_arn
  notification_emails = var.alert_emails
  tags                = local.tags
}

module "codedeploy" {
  source                  = "../../modules/codedeploy"
  name_prefix             = local.name_prefix
  cluster_name            = module.ecs.cluster_name
  service_name            = module.ecs.service_name
  blue_target_group_name  = module.ecs.blue_target_group_name
  green_target_group_name = module.ecs.green_target_group_name
  prod_listener_arn       = module.ecs.prod_listener_arn
  test_listener_arn       = module.ecs.test_listener_arn
  codedeploy_role_arn     = module.iam.codedeploy_role_arn
  tags                    = local.tags
}
