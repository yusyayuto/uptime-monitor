output "sites_table_name" {
  value = aws_dynamodb_table.sites.name
}

output "sites_table_arn" {
  value = aws_dynamodb_table.sites.arn
}

output "status_table_name" {
  value = aws_dynamodb_table.status.name
}

output "status_table_arn" {
  value = aws_dynamodb_table.status.arn
}
