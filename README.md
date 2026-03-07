# OCR API Language PoC

OCR API の実装言語を選定するための PoC リポジトリです。  
Go / TypeScript / Python で同一要件の API を実装し、実装速度・保守性・運用性・実行速度・AWS 連携・型安全性の観点で比較します。

## 目的

本リポジトリの目的は、OCR システムの API 実装言語を選定することです。

OCR ワーカーは非同期で OCR を実行し、API は以下の責務を担当します。

- 署名付き URL の発行
- OCR ジョブ登録
- ジョブステータス取得
- OCR 結果取得

言語選定は印象や好みではなく、同一要件の PoC を通して比較・評価します。

## 前提要件

- クライアントは Web / モバイル
- 画像本体は Cloud Storage 直送
- API は署名付き URL 発行、OCR ジョブ登録、ステータス取得、結果取得を担当
- ワーカーは非同期で OCR 実行
- API は高トラフィックだが重い計算はしない
- 将来的に認証、課金、利用制限、監査ログが必要になる可能性がある

## 比較対象

- Go
- TypeScript (Node.js)
- Python (FastAPI)

## 評価軸

5 段階評価で比較し、以下の重みを設定します。

- 保守性: 15
- 運用性: 15
- 実装速度: 10
- 実行速度: 40
- AWS 連携: 10
- 型安全性: 10

## PoC の対象 API

### 1. 署名付き URL 発行
`POST /v1/uploads/presigned-url`

### 2. OCR ジョブ登録
`POST /v1/ocr-jobs`

### 3. ジョブ状態取得
`GET /v1/ocr-jobs/{jobId}`

### 4. OCR 結果取得
`GET /v1/ocr-jobs/{jobId}/result`

## 実装方針

各言語で以下をできるだけ揃えます。

- 同じ API 仕様
- 同じ DB スキーマ
- 同じレスポンス形式
- 同じエラー形式
- 同じ Docker 前提
- 同じテスト観点

## 想定ディレクトリ構成

```text
.
├── README.md
├── docs/
│   ├── requirements.md
│   ├── evaluation-sheet.md
│   └── api-spec.md
├── go-api/
├── ts-api/
└── py-api/
