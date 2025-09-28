# サイト稼働監視サービス 要件定義および詳細設計

## 概要
- Go言語とAWS ECS Fargateで構築するサイト稼働監視サービス。
- 監視対象のサイト/APIを並行チェックし、異常検知時に通知。
- ゼロダウンタイム更新、スケーラビリティ、可観測性を重視した本番運用想定の最小構成。
- TerraformでインフラをIaC化、GitHub ActionsでCI/CDを整備。

## 機能要件
### ユーザ向け
- 監視対象URL登録: `POST /api/v1/sites` にURL・通知先を登録。
- 監視対象一覧・状態確認: `GET /api/v1/sites`、`GET /api/v1/status/{id}` で最新状態取得。
- 障害通知: 異常検知時にメール通知（SNS経由）。将来的にSlack/Webhook拡張。
- ヘルスチェックAPI: `GET /healthz`, `GET /readyz` を提供。

### 管理者向け
- APIキー管理（発行/無効化）。将来的にユーザ毎のキー、Cognito連携。
- 監視設定統制（登録上限、監視間隔、メンテナンスモード）。
- 全体状況ダッシュボード（監視対象数、成功率、リソース使用率等）。
- WAFルール/シークレット管理（Secrets Managerで安全に保管）。

### 運用担当向け
- サービス稼働監視とSLO/SLA準拠、Runbookに基づくインシデント対応。
- ログ/メトリクス収集: CloudWatch Logs・メトリクスでゴールデンシグナル監視。
- CI/CD監視: GitHub Actions + CodeDeploy のデプロイ状況確認、Blue/Green運用。

## システムアーキテクチャ
- **ネットワーク**: 専用VPC、パブリック/プライベートサブネット。NAT GW経由で外部アクセス。最小権限のセキュリティグループ。
- **エッジ層**: Route53 → CloudFront + WAF。CloudFrontからALBに転送しDDoS/攻撃防御。
- **ロードバランサ**: ALBが入口。CloudFront専用ヘッダ `X-Access-Key` で直接アクセスを遮断。ACM証明書でHTTPS化。
- **アプリ層**: ECS FargateでAPI/ワーカーコンテナを稼働。CloudWatch Logsへ出力。NAT経由で監視対象へアクセス。
- **データ層**: DynamoDBにサイト情報・チェック履歴を格納。オンデマンドキャパシティ、暗号化有効化。
- **通知/可観測性層**: SNSでメール通知。CloudWatchメトリクス・アラーム、ALBアクセスログをS3保存。
- すべてTerraformでコード化し、dev/stg/prod の3環境を用意。

## Terraform構成
- `modules/`: `network`, `ecs`, `alb`, `cloudfront`, `dynamodb`, `sns`, `waf` など機能単位モジュール。
- `envs/{dev,stg,prod}`: 環境ごとにモジュール呼び出し。環境固有の差分（NAT数、DynamoDB容量等）を変数で指定。
- ステートはS3+DynamoDB等のリモートバックエンドで環境別管理。タグ命名規則 `ssk-{env}-{resource}-{detail}`。
- モジュール内は `main.tf`, `variables.tf`, `outputs.tf` 等に分割。機密値はSecrets ManagerやCIの環境変数で注入。

## アプリケーション構成
- Go製。`app/api`, `app/worker`, `app/pkg` の構成を想定。
- APIサーバ機能:
  - `GET /healthz`: Livenessチェック。
  - `GET /readyz`: 依存リソース疎通確認。
  - `POST /api/v1/sites`: 監視対象登録、DynamoDB保存、初期ステータス `pending`。
  - `GET /api/v1/sites`: 登録一覧取得、ページング対応。
  - `GET /api/v1/status/{id}`: 個別最新ステータス取得。
  - 認証: `X-API-Key` ヘッダで静的APIキー検証。Secrets Managerから取得。将来Cognitoで拡張。

## バックグラウンドワーカー
1. 監視対象リストをDynamoDBから取得（60秒間隔）。
2. ゴルーチンで並行HTTPチェック（タイムアウト5秒、非200/エラーで異常判定）。
3. 結果をDynamoDBに保存（バッチ書き込み）。ステータス/最終チェック時刻を更新。
4. 異常時はSNSへPublishしメール通知（将来抑制ロジック追加）。
5. CloudWatch Logsへ構造化ログ、カスタムメトリクスをPut。

## ロギングとエラーハンドリング
- APIは入力バリデーション、HTTPステータス適切化、トレーシングIDでログ紐付け。
- ワーカーは通信/DynamoDB/SNSエラーにリトライ。機密情報はログ出力時にマスク。
- 異常件数もメトリクス化し監視。

## CI/CD
- GitHubリポジトリで `main` ブランチ保護。PRベース開発、Issue/PRテンプレート整備。
- GitHub Actions CI: Goテスト・lint→マルチステージDockerビルド→ECRへプッシュ。OIDC+IAMロールで認証。
- デリバリ: CodeDeployを利用したECS Blue/Green。
  - 新タスク定義登録→デプロイ開始→`/readyz` ヘルスチェック→トラフィック切替→旧タスク停止。
  - 失敗時自動ロールバック、通知。
  - `stg` → `prod` の順でデプロイ、環境ごとにSecrets/変数を切替。
- 将来: 統合/負荷テスト、Terraform Plan承認、コンテナ脆弱性スキャン等を組み込み。

## IAM・セキュリティ・運用
- 最小権限: ECSタスクロールにDynamoDB/SNS/Logsの必要権限のみ。Terraform実行ロールはMFA必須。
- GitHub Actions OIDCロールはECR/CodeDeploy権限に限定し信頼ポリシー厳格化。Access Analyzerで過剰権限検知。
- 多層防御: CloudFront+WAF、セキュリティグループ、APIキー認証。すべてHTTPS、Secrets Managerでシークレット管理。
- DynamoDB暗号化、ログの個人情報最小化。ACMでTLS1.2以上。CloudTrail/AWS Configで監査。
- CloudWatchアラーム: API可用性、通知遅延、連続失敗数など。Dead Man's Switch的外形監視も検討。
- Runbook整備、ゲームデー実施。脆弱性対応やTerraformモジュール更新を継続。

## 将来的な拡張
- 通知チャネル追加（Slack, Webhook, PagerDuty等）。
- Cognito導入でユーザ管理強化、Webダッシュボード提供。
- EventBridge+SQSによる分散処理、マルチリージョンDR検討。
- 監視データ分析基盤構築、機械学習異常検知。
- コスト最適化（Lambda化、オンデマンドジョブ）。

