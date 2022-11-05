CREATE TABLE "users" (
    id VARCHAR primary key,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    access_token VARCHAR,
    refresh_token VARCHAR,
    strava_id VARCHAR,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);