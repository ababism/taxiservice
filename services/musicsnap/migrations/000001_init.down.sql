DROP INDEX IF EXISTS idx_event_authors_event_id;
DROP INDEX IF EXISTS idx_reviews_user_id;
DROP INDEX IF EXISTS idx_subscriptions_subscriber_id;
DROP INDEX IF EXISTS idx_playlist_items_playlist_id;
DROP INDEX IF EXISTS idx_reviews_piece_id;
DROP INDEX IF EXISTS idx_ratings_user_id;

-- DROP INDEX IF EXISTS idx_users_email;
-- DROP INDEX IF EXISTS idx_users_username;
-- DROP INDEX IF EXISTS idx_ratings_piece_id;
-- DROP INDEX IF EXISTS idx_reviews_description_id;

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS event_photos;
DROP TABLE IF EXISTS event_authors;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS playlist_items;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS reactions;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS descriptions;
DROP TABLE IF EXISTS ratings;
-- DROP TABLE IF EXISTS pieces;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;




