CREATE TABLE IF NOT EXISTS Users
(
    user_id  SERIAL PRIMARY KEY,
    username VARCHAR(50)  NOT NULL,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Messages
(
    message_id   SERIAL PRIMARY KEY,
    sender_id    INTEGER NOT NULL,
    receiver_id  INTEGER NOT NULL,
    message_text TEXT    NOT NULL,
    sent_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    readed       BOOLEAN   DEFAULT FALSE,
    FOREIGN KEY (sender_id) REFERENCES Users (user_id),
    FOREIGN KEY (receiver_id) REFERENCES Users (user_id)
);