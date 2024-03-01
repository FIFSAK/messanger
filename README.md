# Messanger project
## Rest API
- **Register a New User**
  - `POST /register`
- **Login**
  - `GET /login`
- **Update User Authentication Type**
  - `PATCH /login/{type}`
- **Delete User**
  - `DELETE /login`

### User Information
- **Retrieve All Users**
  - `GET /users`

### Messaging
- **Send a Message**
  - `POST /message/{id}`
- **Update a Message**
  - `UPDATE /message/{id}`
- **Delete a Message**
  - `DELETE /message/{id}`

### Notifications
- **Retrieve Notifications**
  - `GET /notifications`
Database Structure
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
  sent_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
}
