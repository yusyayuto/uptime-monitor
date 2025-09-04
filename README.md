# uptime-monitor — Go × AWS (ECS/Fargate) サイト稼働監視サービス

複数サイトを並行チェックし、異常を検知すると通知するAPI & Worker。
- インフラ: CloudFront + WAF + ALB + ECS(Fargate)
- IaC: Terraform（3環境: dev / stg / prod）
- CI/CD: GitHub Actions → CodeDeploy (Blue/Green)
- データ: Auroraで開始、将来は DynamoDB（時系列最適化）へ
- セキュリティ: ALB×Cognito(OAuth2.1), 直叩き防止, 最小権限IAM

## 進め方（フェーズ）
0. 準備（このREADME・雛形・命名/タグ規約）
1. VPC/CF/WAF/ALB（/healthz到達）
2. ECS Hello + 直叩き防止
3. CI/CD（BGデプロイ）
4. アプリMVP（Aurora, SNS通知）
5. SLO/ダッシュボード/アラーム
6. DynamoDB移行（時系列最適化）
7. EventBridge+SQS（キュー深度スケール）
8. 3環境運用 & 最終ドキュメント

## ディレクトリ
- `docs/` : 図・設計（architecture, naming, slo, runbook）
- `infra/` : Terraform（modules と envs/dev|stg|prod）
- `app/`   : Go（api, worker, pkg 共通）
- `.github/` : CI/CD と Issue/PR テンプレ

## 命名・タグ（詳しくは docs/naming.md）
- 命名: `ssk-{env}-{resource}-{detail}`（env: dev/stg/prod）
- タグ: `System=ssk, Service=uptime-monitor, Env=*, Owner=ssk`
