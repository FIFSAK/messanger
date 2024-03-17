# Messanger project

## Rest API
- `GET /health-check`
- `POST /register`
- `GET /login`
- `PATCH /login/{type}`
- `DELETE /login`
- `GET /users`
- `GET /message/send`
- `GET /message/received`
- `POST /message`
- `PATCH /message`
- `DELETE /message`
- `GET /message/notifications`

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

## api documentation
- **GET /health-check**
```json lines
OK
```
- **POST /register**
```json lines
{
  "login": "your_username",
  "password": "your_password"
}
```
- **GET /login** returns jwt token
```json lines
{
  "login": "your_username",
  "password": "your_password"
}
```
- **PATCH /login/{type}** type can be `login` or `password` if you want update login or password respectively
```json lines
header {
"Bearer <Your API key>"
}
{
  "new-login": "new_username",
  "new-password": "new_password"
}
```
- **DELETE /login**
```json lines
header {
    "Bearer <Your API key>"
}
```
- **GET /users**
```json lines
header {
    "Bearer <Your API key>"
}
```
- **GET /message/send** return all your sent messages
```json lines 
header {
    "Bearer <Your API key>"
}
```
- **GET /message/received** return all your received messages
```json lines
header {
    "Bearer <Your API key>"
}
```
- **POST /message**
```json lines
header {
    "Bearer <Your API key>"
}
{
  "receiver_id": "receiver_id",
  "message_text": "your_message"
}
```
- **PATCH /message**
```json lines
header {
    "Bearer <Your API key>"
}
{
  "message_id": "message_id",
  "message_text": "your_message"
}
```
- **DELETE /message**
```json lines
header {
    "Bearer <Your API key>"
}
{
  "message_id": "message_id"
}
```
- **GET /message/notifications** return all your unread messages
```json lines
header {
    "Bearer <Your API key>"
}
```









