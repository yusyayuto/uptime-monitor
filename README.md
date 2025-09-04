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
