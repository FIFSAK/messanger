# Messenger Project API Documentation

## REST API Endpoints

### Health Check
- `GET /health-check`
    - **Response:** `OK`

### User Registration
- `POST /register`
    - **Body:** 
```json lines 
{ 
  "login": "your_username",
  "password": "your_password"
}
 ```

### User Login
- `GET /login`
    - **Description:** Returns a JWT token.
    - **Body:** 
```json lines 
{ 
  "login": "your_username",
  "password": "your_password"
}
 ```
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
    - **Body:**
```json lines 
{
"new-login": "your_username",
"new-password": "your_password"
}
 ```

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
    - **Body:**
```json lines 
{ 
  "receiver_id": "receiver_id",
  "message_text": "your_message"
}
 ```

### Update a Message
- `PATCH /message`
    - **Header:** `Authorization: Bearer <Your API key>`
    - **Body:**
```json lines 
{ 
  "message_id": "message_id",
  "message_text": "your_message"
}
 ```
### Delete a Message
- `DELETE /message`
    - **Header:** `Authorization: Bearer <Your API key>`
    - **Body:**
```json lines 
{ 
  "message_id": "message_id",
}
 ```

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
