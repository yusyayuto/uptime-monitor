resource "aws_dynamodb_table" "sites" {
  name         = "${var.name_prefix}-sites"
  billing_mode = var.billing_mode
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-sites"
  })
}

resource "aws_dynamodb_table" "status" {
  name         = "${var.name_prefix}-status"
  billing_mode = var.billing_mode
  hash_key     = "site_id"

  attribute {
    name = "site_id"
    type = "S"
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-status"
  })
}
