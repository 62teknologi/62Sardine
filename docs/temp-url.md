## Temporary URL

- **URL:** `/api/v1/temporary-url`
- **Method:** `GET`

### Params

| Name                | Type    | Required | Description                                                           |
| ------------------- | ------- | -------- | --------------------------------------------------------------------- |
| `path`              | String  | Yes      | The path to the targeted file.                                        |
| `expired_in_minute` | Integer | No       | The expiration time for the uploaded file in minutes. Default is `30` |

### Example

#### ENV

```env
HTTP_SERVER_ADDRESS="0.0.0.0:10082"
FILESYSTEM_DISK="local"
FILESYSTEM_FOLDER=sardine-test
EXPORT_FOLDER=storage
APP_URL=http://localhost:8000
```

#### Request

```http
GET /api/v1/temporary-url?path=sardine-test/test-folder/example.jpg&expired_in_minute=30
```

#### Response

```json
{
  "data": {
    "url": "http://localhost:8000/storage/app/sardine-test/test-folder/example.jpg"
  }
}
```
