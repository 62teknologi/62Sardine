## Upload

- **URL:** `/api/v1/files`
- **Method:** `POST`
- **Content-Type:** `multipart/form-data`

### Body

| Name            | Type    | Required | Description                                                                                     |
| --------------- | ------- | -------- | ----------------------------------------------------------------------------------------------- |
| `file`          | File    | Yes      | The file to be uploaded.                                                                        |
| `visibility`    | String  | Yes      | The visibility setting for the uploaded file (e.g., "public", "private"). Default is "private". |
| `folder`        | String  | No       | The folder to which the file belongs.                                                           |
| `file_name`     | String  | No       | The desired filename for the uploaded file.                                                     |
| `resize_width`  | Integer | No       | The width to which the uploaded image should be resized.                                        |
| `resize_height` | Integer | No       | The height to which the uploaded image should be resized.                                       |
| `compress`      | Integer | No       | The desired byte size for compression of the uploaded file.                                     |
| `accept`        | String  | No       | The accepted file type of the uploaded file.                                                   |

### Example

#### ENV

```env
HTTP_SERVER_ADDRESS="0.0.0.0:10082"
FILESYSTEM_DISK="local"
FILESYSTEM_FOLDER=sardine-test
EXPORT_FOLDER=storage
APP_URL=http://localhost:8000
```

#### Request Body

```curl
curl --location 'http://localhost:10082/api/v1/files' \
--header 'Content-Type: multipart/form-data' \
--form 'file=@"example.jpeg"' \
--form 'visibility="public"' \
--form 'folder="test-folder"' \
--form 'file_name="my_custom_name"'
```

#### Response

```json
{
  "data": {
    "bucket": "",
    "client_original_extention": "jpeg",
    "client_original_name": "example.jpeg",
    "content_type": "image/jpeg",
    "disk": "local",
    "extension": "jpg",
    "file_name": "my_custom_name",
    "more_info": {
      "height": 1350,
      "width": 1080
    },
    "path": "sardine-test/test-folder/my_custom_name.jpg",
    "size": 279643,
    "url": "http://localhost:8000/storage/app/sardine-test/test-folder/my_custom_name.jpg"
  }
}
```
