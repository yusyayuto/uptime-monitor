# 命名規約（Classmethod流に準拠・簡略版）

## フォーマット
`ssk-{env}-{resource}-{detail}`
- env: `dev` / `stg` / `prod`
- resource: `vpc|alb|ecs|rds|cf|waf|s3|iam|ecr|sg|dynamodb|sqs|eb|sm` など
- detail: 任意の識別（用途/サブ名）

### 例
- `ssk-prod-vpc`
- `ssk-stg-alb-api`
- `ssk-dev-ecs-api-service`
- `ssk-prod-dynamodb-checks`
- `ssk-prod-s3-alb-logs`

## タグ共通化
- `System = ssk`
- `Service = uptime-monitor`
- `Env = dev|stg|prod`
- `Owner = ssk`
