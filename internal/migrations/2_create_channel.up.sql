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
    read       BOOLEAN   DEFAULT FALSE,
    FOREIGN KEY (sender_id) REFERENCES Users (user_id),
    FOREIGN KEY (receiver_id) REFERENCES Users (user_id)
);

CREATE TABLE IF NOT EXISTS Channel
(
    chanel_id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    chanel_name VARCHAR(50) UNIQUE NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES Users(user_id)
);


CREATE TABLE IF NOT EXISTS ChannelUsers
(
    chanel_user_id SERIAL PRIMARY KEY,
    chanel_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (chanel_id) REFERENCES Channel(chanel_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);


CREATE TABLE IF NOT EXISTS ChannelMessages
(
    chanel_message_id SERIAL PRIMARY KEY,
    chanel_id INTEGER NOT NULL,
    message_text TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chanel_id) REFERENCES Channel(chanel_id)
);