## Delete

- **URL:** `/api/v1/files`
- **Method:** `DELETE`

### Params

| Name   | Type   | Required | Description                    |
| ------ | ------ | -------- | ------------------------------ |
| `path` | String | Yes      | The path to the targeted file. |

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
GET /api/v1/files?path=sardine-test/test-folder/my_custom_name.jpg
```

#### Response

```json
{
  "data": {
    "success": true
  }
}
```
