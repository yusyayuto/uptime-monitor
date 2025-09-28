output "cluster_name" {
  value = aws_ecs_cluster.this.name
}

output "service_name" {
  value = aws_ecs_service.api.name
}

output "blue_target_group_arn" {
  value = aws_lb_target_group.blue.arn
}

output "blue_target_group_name" {
  value = aws_lb_target_group.blue.name
}

output "green_target_group_arn" {
  value = aws_lb_target_group.green.arn
}

output "green_target_group_name" {
  value = aws_lb_target_group.green.name
}

output "load_balancer_arn" {
  value = aws_lb.api.arn
}

output "load_balancer_dns_name" {
  value = aws_lb.api.dns_name
}

output "prod_listener_arn" {
  value = aws_lb_listener.prod.arn
}

output "test_listener_arn" {
  value = aws_lb_listener.test.arn
}
