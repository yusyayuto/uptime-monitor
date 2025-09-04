# uptime-monitor — Go × AWS (ECS/Fargate) サイト稼働監視サービス

> **uptime-monitor** は、複数サイトを並行チェックして異常を検知・通知する **Go × AWS(ECS/Fargate)** の監視サービスです。  
> ポートフォリオ目的で「ゼロダウン更新・スケール・可観測性」を重視した **本番想定の最小構成** を採用します。

---

## Quickstart (予定)

> 実装が入り次第、コマンドを追記します。まずは雛形だけ先置き。

```bash
# APIローカル起動（後で実装）
make run-api

# ワーカー（後で実装）
make run-worker

# ヘルスチェック
curl -s http://localhost:8080/healthz
```

## Status / Roadmap

* **フェーズ0**: リポ雛形・命名規約・Issue/PRテンプレ
* **フェーズ1**: VPC / CloudFront / WAF / ALB（/healthz到達）
* **フェーズ2**: ECS(Fargate) Hello + 直叩き防止
* **フェーズ3**: CI/CD (GitHub Actions → CodeDeploy Blue/Green)
* **フェーズ4**: アプリMVP（Aurora 版 / SNS 通知）
* **フェーズ5**: SLO・ダッシュボード・アラーム
* **フェーズ6**: DynamoDB移行（時系列最適化）
* **フェーズ7**: EventBridge + SQS（キュー深度スケール）
* **フェーズ8**: 3環境運用・最終ドキュメント

## Tech Decisions（要点）

* **Go**: 軽量単一バイナリ＆並行処理に強い（多対象の同時監視に最適）
* **ECS Fargate**: サーバ管理不要。本番でのスケール/Blue-Greenが容易
* **Terraform**: IaCで再現性・レビュー性を担保（3環境 dev/stg/prod）
* **Aurora → DynamoDB**: まずRDBで開発容易性 → 将来は時系列最適化へ移行
* **CloudFront + WAF + ALB**: 境界で防御。ALB直叩きはカスタムヘッダで拒否
* **GitHub Actions + CodeDeploy**: OIDC AssumeRoleで安全にBGデプロイ

## API (MVP 予定)

* `GET /healthz` : 常に200（Liveness）
* `GET /readyz` : 依存先OKで200（Blue/Green判定）
* `POST /api/v1/sites` : 監視対象登録
* `GET /api/v1/sites` : 一覧
* `GET /api/v1/status/{id}` : 直近ステータス

## Environments / Naming / Tags

* **環境**: `dev` / `stg` / `prod`
* **命名**: `ssk-{env}-{resource}-{detail}`
  * 例: `ssk-stg-alb-api`, `ssk-prod-dynamodb-checks`
* **タグ**: `System=ssk, Service=uptime-monitor, Env=*, Owner=ssk`

## CI/CD

`git push main` → GitHub Actions (test/build) → ECR push → CodeDeploy (ECS Blue/Green, 443/8443, 自動ロールバック)

## Security（方針）

* CloudFront + WAF → ALB（**カスタムヘッダ一致時のみ許可**）
* ALB × Cognito (OAuth2.1) で `/api/*` を保護
* IAM 最小権限・Secrets Manager でDB資格情報管理

## SLO（自己監視・予定）

* API 可用性: **99.9%**
* 取込成功率: **99.5%**
* 通知遅延: **P95 < 30s**
* 近似バーンレート: **1h > 2% / 6h > 5% / 3d > 10%**

詳細は `docs/slo.md` を参照（追記予定）

## Non-Goals / Cost

* 個人運用のため **NATは1台/最小構成**。高額な常時大規模負荷試験は非対象
* mTLS / 証明書ピニング / AI予測は **将来の拡張**（設計メモのみ）

## Dev（予定）

* `docker-compose` で MySQL 起動
* `make migrate` で DB 初期化
* `.env` で API ポート・DB 接続設定

## Repository Structure

```
docs/            # 図・仕様・設計（architecture, naming, slo, runbook）
infra/           # Terraform（modules と envs/dev|stg|prod）
app/             # Go（api, worker, pkg 共通）
.github/         # CI/CD と Issue/PR テンプレ
```

## License

MIT

---

# 📦 反映のやり方（コマンド例）

```bash
# 1) ローカルでブランチを切る（任意）
git checkout -b chore/update-readme

# 2) README.md を上の内容で丸ごと置き換え → 保存

# 3) コミット＆プッシュ
git add README.md
git commit -m "docs: bootstrap README with roadmap & decisions"
git push -u origin chore/update-readme

# 4) GitHub で Pull Request 作成 → マージ
