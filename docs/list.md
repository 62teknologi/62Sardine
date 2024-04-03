## List

- **URL:** `/api/v1/files`
- **Method:** `GET`

### Params

| Name   | Type   | Required | Description                      |
| ------ | ------ | -------- | -------------------------------- |
| `path` | String | Yes      | The path to the targeted folder. |

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
GET /api/v1/files?path=sardine-test/test-folder
```

#### Response

```json
{
  "data": ["my_custom_name.jpg"]
}
```
