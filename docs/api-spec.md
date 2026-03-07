# API Specification

OCR API PoC 用の API 仕様です。

この API は OCR 処理のジョブ管理を目的としています。  
OCR 本体はワーカーが実行し、API はジョブ管理と結果取得を担当します。

---

# 共通レスポンス形式

## Success

```json
{
  "data": {}
}
```

## Error

```json
{
  "error": {
    "code": "invalid_request",
    "message": "filename is required"
  }
}
```

---

# 1. Presigned URL 発行

アップロード用の署名付き URL を発行します。

## Endpoint

POST /v1/uploads/presigned-url

## Request

```json
{
  "filename": "receipt.jpg",
  "contentType": "image/jpeg"
}
```

## Response

```json
{
  "data": {
    "objectKey": "uploads/2026/03/07/uuid-receipt.jpg",
    "uploadUrl": "https://storage.example.com/....",
    "expiresIn": 300
  }
}
```

## Description

クライアントはこの URL を使用して Cloud Storage に直接ファイルをアップロードします。

---

# 2. OCRジョブ登録

アップロード済み画像をもとに OCR ジョブを作成します。

## Endpoint

POST /v1/ocr-jobs

## Request

```json
{
  "objectKey": "uploads/2026/03/07/uuid-receipt.jpg"
}
```

## Response

```json
{
  "data": {
    "jobId": "job_xxxxx",
    "status": "queued"
  }
}
```

## Description

この API は以下を行います。

1. OCRジョブをDBに登録
2. Queue にジョブを送信
3. jobId を返す

---

# 3. ジョブステータス取得

OCR ジョブの状態を取得します。

## Endpoint

GET /v1/ocr-jobs/{jobId}

## Response

```json
{
  "data": {
    "jobId": "job_xxxxx",
    "status": "processing"
  }
}
```

---

# 4. OCR結果取得

OCR 処理結果を取得します。

## Endpoint

GET /v1/ocr-jobs/{jobId}/result

## Response

```json
{
  "data": {
    "jobId": "job_xxxxx",
    "status": "succeeded",
    "result": {
      "text": "sample OCR result"
    }
  }
}
```

---

# Job Status

| status | description |
|------|-------------|
| queued | ジョブ作成済み |
| processing | OCR処理中 |
| succeeded | OCR完了 |
| failed | OCR失敗 |

---

# Database Schema

テーブル: `ocr_jobs`

| column | type | description |
|------|------|-------------|
| id | string | ジョブID |
| object_key | string | Cloud Storage オブジェクトキー |
| status | string | ジョブ状態 |
| result_json | json | OCR結果 |
| created_at | timestamp | 作成日時 |
| updated_at | timestamp | 更新日時 |

---

# 処理フロー

1. クライアントが Presigned URL を取得
2. クライアントが Cloud Storage に画像アップロード
3. クライアントが OCR ジョブ作成 API を呼ぶ
4. API が Queue にジョブを送信
5. Worker が OCR を実行
6. Worker が結果を DB に保存
7. クライアントが結果取得 API を呼ぶ
