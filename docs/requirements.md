# Requirements

本ドキュメントは OCR API 言語選定 PoC の要件を整理するためのものです。

この PoC の目的は、API 実装言語を比較し、最適な技術選定を行うことです。

比較対象言語

- Go
- TypeScript (Node.js)
- Python (FastAPI)

---

# システム概要

OCR システムは以下の構成を想定しています。

Client (Web / Mobile)
        │
        ▼
       API
        │
        ├── Cloud Storage (画像アップロード)
        │
        ├── Queue (OCR ジョブ送信)
        │
        ▼
     OCR Worker
        │
        ▼
       Database

---

# APIの役割

API は以下の責務を持ちます。

- 署名付きアップロード URL 発行
- OCR ジョブ登録
- ジョブステータス取得
- OCR 結果取得

OCR 処理そのものはワーカーが実行します。

---

# 非機能要件

## 高トラフィック対応

API は高トラフィックを想定しますが、CPU負荷の高い処理は行いません。

主な処理

- DBアクセス
- Storage連携
- Queue送信

---

## 拡張性

将来的に以下の機能を追加する可能性があります。

- 認証
- APIキー
- 課金
- 利用制限
- 監査ログ

そのため、APIは保守しやすい構造が求められます。

---

# 必須条件

PoC では以下を満たすことを条件とします。

- PostgreSQL と連携できる
- Cloud Storage と連携できる
- Queue にジョブ送信できる
- OpenAPI を管理できる
- 入力バリデーションが実装できる
- テストが書ける
- 構造化ログが出力できる
- Docker で起動できる

---

# 評価対象

PoC では以下を評価対象とします。

- API実装の容易さ
- バリデーション実装
- OpenAPI 管理
- DBアクセス
- Storage連携
- Queue送信
- ログ設計
- テストの書きやすさ
- Docker運用

---

# 対象外

PoCでは以下は対象外とします。

- OCRエンジン
- AIモデル
- 認証
- 課金
- 本番監視
- CI/CD

これらは本実装フェーズで検討します。
