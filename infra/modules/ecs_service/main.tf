resource "aws_lb" "api" {
  name               = "${var.name_prefix}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-alb"
  })
}

resource "aws_lb_target_group" "blue" {
  name        = "${var.name_prefix}-blue"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    healthy_threshold   = 3
    unhealthy_threshold = 2
    timeout             = 5
    path                = "/readyz"
    matcher             = "200"
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-blue"
  })
}

resource "aws_lb_target_group" "green" {
  name        = "${var.name_prefix}-green"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    healthy_threshold   = 3
    unhealthy_threshold = 2
    timeout             = 5
    path                = "/readyz"
    matcher             = "200"
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-green"
  })
}

resource "aws_lb_listener" "prod" {
  load_balancer_arn = aws_lb.api.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.blue.arn
  }
}

resource "aws_lb_listener" "test" {
  load_balancer_arn = aws_lb.api.arn
  port              = var.test_listener_port
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.green.arn
  }
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "/aws/ecs/${var.name_prefix}-api"
  retention_in_days = 30

  tags = var.tags
}

resource "aws_ecs_cluster" "this" {
  name = "${var.name_prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = var.tags
}

locals {
  container_definitions = [
    {
      name      = "api"
      image     = var.container_image
      essential = true
      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]
      environment = [
        { name = "APP_ENV", value = var.tags["Env"] },
        { name = "API_PORT", value = "8080" },
        { name = "SITES_TABLE", value = "${var.name_prefix}-sites" },
        { name = "STATUS_TABLE", value = "${var.name_prefix}-status" }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.api.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "api"
        }
      }
    },
    {
      name      = "worker"
      image     = var.container_image
      essential = true
      command   = ["/app", "worker"]
      environment = [
        { name = "APP_ENV", value = var.tags["Env"] },
        { name = "POLL_INTERVAL_SECONDS", value = "60" },
        { name = "SITES_TABLE", value = "${var.name_prefix}-sites" },
        { name = "STATUS_TABLE", value = "${var.name_prefix}-status" }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.api.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "worker"
        }
      }
    }
  ]
}

data "aws_region" "current" {}

resource "aws_ecs_task_definition" "api" {
  family                   = "${var.name_prefix}-task"
  cpu                      = var.cpu
  memory                   = var.memory
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = var.task_execution_role_arn
  task_role_arn            = var.task_role_arn

  container_definitions = jsonencode(local.container_definitions)

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  tags = var.tags
}

resource "aws_ecs_service" "api" {
  name            = "${var.name_prefix}-svc"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    assign_public_ip = false
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_security_group_id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.blue.arn
    container_name   = "api"
    container_port   = 8080
  }

  deployment_controller {
    type = "CODE_DEPLOY"
  }

  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = var.tags
}
