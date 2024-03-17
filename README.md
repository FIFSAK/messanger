# Messenger Project API Documentation

## REST API Endpoints

### Health Check
- `GET /health-check`
    - **Response:** `OK`

### User Registration
- `POST /register`    
  - **Body:** Use `multipart/form-data` to send the following fields:
    - `login` (text): your username
    - `password` (text): your password

### User Login
- `GET /login`
    - **Description:** Returns a JWT token.
    - **Body:** Use `multipart/form-data` to send the following fields:
      - `login` (text): your username
      - `password` (text): your password
- **Response:**
```json lines
{
  "token": "your_jwt_token"
}
```

### Update Login or Password
- `PATCH /login/{type}`
    - **Description:** `{type}` can be `login` or `password`.
    - **Header:** `Authorization: Bearer <Your API key>`
    - **Body:** Use `multipart/form-data` to send the following fields:
      - `new-login` (text): your new login
      - `new-password` (text): your new password

### Delete Login
- `DELETE /login`
    - **Header:** `Authorization: Bearer <Your API key>`

### List Users
- `GET /users`
    - **Header:** `Authorization: Bearer <Your API key>`

### Sent Messages
- `GET /message/send`
    - **Header:** `Authorization: Bearer <Your API key>`

### Received Messages
- `GET /message/received`
    - **Header:** `Authorization: Bearer <Your API key>`

### Send a Message
- `POST /message`
  - **Header:** `Authorization: Bearer <Your API key>`
  - **Body:** Use `multipart/form-data` to send the following fields:
    - `receiver_id` (text): ID of the receiver
    - `message_text` (text): The message content

### Update a Message
- `PATCH /message`
    - **Header:** `Authorization: Bearer <Your API key>`
    - **Body:** Use `multipart/form-data` to send the following fields:
      - `message_id` (text): ID of the message
      - `message_text` (text): The message content

### Delete a Message
- `DELETE /message`
    - **Header:** `Authorization: Bearer <Your API key>`
    - **Body:**Use `multipart/form-data` to send the following fields:
      - `message_id` (text): ID of the message

### Unread Messages Notifications
- `GET /message/notifications`
    - **Header:** `Authorization: Bearer <Your API key>`
    - **Response:** list of unread messages


## Database Structure

```sql
Table Users {
  user_id SERIAL [pk]
  username VARCHAR(50) [not null]
  password VARCHAR(100) [not null]
}

Table Messages {
  message_id SERIAL [pk]
  sender_id INTEGER [not null, ref: > Users.user_id]
  receiver_id INTEGER [not null, ref: > Users.user_id]
  message_text TEXT [not null]
  readed BOOLEAN [default: false]
  sent_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
}
```
## Environment variables
**.env** in the root folder of the project
````
host=fullstack-postgres
dbname=your_db_name
sslmode=disable
port=5432
user=your_user
password=your_password
secret=your_256_bit_secret
PGADMIN_DEFAULT_EMAIL=your_email
PGADMIN_DEFAULT_PASSWORD=your_password
````
## Run project

**Start project first time or after changes** ``` docker-compose up --build```

**otherwise** ```docker-compose up```

