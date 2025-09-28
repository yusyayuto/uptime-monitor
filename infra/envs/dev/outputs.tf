output "vpc_id" {
  value = module.networking.vpc_id
}

output "alb_dns_name" {
  value = module.ecs.load_balancer_dns_name
}

output "sns_topic_arn" {
  value = module.observability.sns_topic_arn
}

output "codedeploy_application" {
  value = module.codedeploy.application_name
}

output "codedeploy_deployment_group" {
  value = module.codedeploy.deployment_group_name
}
