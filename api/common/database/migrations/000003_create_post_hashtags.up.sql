CREATE TABLE post_hashtags (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id    uuid        NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    tag        text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_post_hashtags_tag ON post_hashtags (tag);
CREATE INDEX idx_post_hashtags_post_id ON post_hashtags (post_id);
