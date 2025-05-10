-- Пользователи
CREATE TABLE users
(
    id             UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    nickname       TEXT      NOT NULL UNIQUE,
    avatar_url     VARCHAR(255),
    background_url VARCHAR(255),
--     avatar_url     TEXT,
--     background_url TEXT,
    bio            TEXT,

    email          TEXT      NOT NULL UNIQUE,
    password_hash  TEXT      NOT NULL,
    created_at     TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

-- user roles table
CREATE TABLE user_roles
(
    user_id    UUID REFERENCES users (id),
    role       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    PRIMARY KEY (user_id, role)
);


-- Оценки
CREATE TABLE ratings
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES users (id),
    piece_id   VARCHAR(255) NOT NULL,
    rating     INTEGER      NOT NULL CHECK (rating BETWEEN 1 AND 10),
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
-- Описания
CREATE TABLE descriptions
(
    id         UUID PRIMARY KEY,
    text       TEXT      NOT NULL,
    photo_link TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- Треды
CREATE TABLE threads
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Рецензии
CREATE TABLE reviews
(
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id    UUID REFERENCES users (id),
    piece_id   VARCHAR(255) NOT NULL,

    rating     INTEGER      NOT NULL CHECK (rating BETWEEN 1 AND 10),
    photo_url  VARCHAR(255),
    content    TEXT,

    moderated  BOOLEAN      NOT NULL DEFAULT false,
    published  BOOLEAN      NOT NULL DEFAULT true,

--     description_id UUID REFERENCES descriptions (id),
--     thread_id      UUID REFERENCES threads (id),
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),

    CONSTRAINT content_length CHECK (char_length(content) <= 10000)
);


-- Комментарии
CREATE TABLE comments
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES users (id),
    thread_id  UUID REFERENCES threads (id),
    text       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Реакции
CREATE TABLE reactions
(
    id         INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id    UUID REFERENCES users (id),
    review_id  INT REFERENCES reviews (id),
    type       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT reaction_type CHECK (type IN ('like', 'dislike')),
    CONSTRAINT unique_reaction UNIQUE (user_id, review_id)
);

-- Подписки
CREATE TABLE subscriptions
(
--     sub_id                SERIAL PRIMARY KEY,
    sub_id            INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    subscriber_id     UUID REFERENCES users (id),
    followed_id       UUID REFERENCES users (id),
    notification_flag BOOLEAN   NOT NULL DEFAULT false,
    created_at        TIMESTAMP NOT NULL,
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_subscription UNIQUE (subscriber_id, followed_id)
);

-- Плейлисты
CREATE TABLE playlists
(
    id          UUID PRIMARY KEY,
    user_id     UUID REFERENCES users (id),
    name        TEXT      NOT NULL,
    description TEXT,
    cover_url   TEXT,
    is_ranked   BOOLEAN   NOT NULL,
    is_private  BOOLEAN   NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Элементы плейлиста
CREATE TABLE playlist_items
(
    id             UUID PRIMARY KEY,
    playlist_id    UUID REFERENCES playlists (id),
    piece_id       VARCHAR(255) NOT NULL,
    description_id UUID REFERENCES descriptions (id),
    type           TEXT         NOT NULL CHECK (type IN ('Piece', 'Description')),
    created_at     TIMESTAMP    NOT NULL,
    updated_at     TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- События
CREATE TABLE events
(
    id          UUID PRIMARY KEY,
    name        TEXT      NOT NULL,
    date        TIMESTAMP NOT NULL,
    cover_url   TEXT,
    ticket_link TEXT,
    map_link    TEXT,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Авторы событий
CREATE TABLE event_authors
(
    id         UUID PRIMARY KEY,
    event_id   UUID REFERENCES events (id),
    user_id    UUID REFERENCES users (id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Фото событий
CREATE TABLE event_photos
(
    id         UUID PRIMARY KEY,
    event_id   UUID REFERENCES events (id),
    photo_url  TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Уведомления
CREATE TABLE notifications
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES users (id),
    type       TEXT      NOT NULL,
    message    TEXT      NOT NULL,
    read       BOOLEAN   NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- -- Треки
-- CREATE TABLE tracks (
--                         id UUID PRIMARY KEY REFERENCES pieces(id),
--                         title TEXT NOT NULL,
--                         duration_ms INTEGER NOT NULL,
--                         explicit BOOLEAN NOT NULL,
--                         created_at TIMESTAMP NOT NULL,
--                         updated_at TIMESTAMP NOT NULL
-- );
--
-- -- Альбомы
-- CREATE TABLE albums (
--                         id UUID PRIMARY KEY REFERENCES pieces(id),
--                         title TEXT NOT NULL,
--                         cover_url TEXT,
--                         release_date TIMESTAMP NOT NULL,
--                         created_at TIMESTAMP NOT NULL,
--                         updated_at TIMESTAMP NOT NULL
-- );

-- Индексы для ускорения поиска
CREATE INDEX idx_ratings_user_id ON ratings (user_id);
CREATE INDEX idx_reviews_piece_id ON reviews (piece_id);
CREATE INDEX idx_playlist_items_playlist_id ON playlist_items (playlist_id);
CREATE INDEX idx_subscriptions_subscriber_id ON subscriptions (subscriber_id);
CREATE INDEX idx_reviews_user_id ON reviews (user_id);
CREATE INDEX idx_event_authors_event_id ON event_authors (event_id);


CREATE
    OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
        = NOW();
    RETURN NEW;
END ;
$$
    language 'plpgsql';

-- крутой метод
DO
$$
    DECLARE
        cur_table_name TEXT;
    BEGIN
        FOR cur_table_name IN
            SELECT table_name
            FROM information_schema.columns
            WHERE column_name = 'updated_at'
              AND table_schema = 'public'
            LOOP
                EXECUTE format(
                        'CREATE TRIGGER update_%I_updated_at
                         BEFORE UPDATE ON %I
                         FOR EACH ROW
                         EXECUTE FUNCTION update_updated_at();',
                        cur_table_name, cur_table_name
                        );
            END LOOP;
    END;
$$;


CREATE
    OR REPLACE FUNCTION update_created_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.created_at
        = NOW();
    RETURN NEW;
END;
$$
    language 'plpgsql';

--  write ctreate trigger for created_at
DO
$$
    DECLARE
        cur_table_name TEXT;
    BEGIN
        FOR cur_table_name IN
            SELECT table_name
            FROM information_schema.columns
            WHERE column_name = 'created_at'
              AND table_schema = 'public'
            LOOP
                EXECUTE format(
                        'CREATE TRIGGER create_%I_created_at
                         BEFORE INSERT ON %I
                         FOR EACH ROW
                         EXECUTE FUNCTION update_created_at();',
                        cur_table_name, cur_table_name
                        );
            END LOOP;
    END;
$$;
