package database

var ReactionTable = `
    CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    target_type TEXT NOT NULL CHECK(target_type IN ('post', 'comment')),
    reaction_type TEXT NOT NULL CHECK(reaction_type IN ('like', 'dislike')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    UNIQUE (user_id, post_id, target_type),
    UNIQUE (user_id, comment_id, target_type),
    CHECK (
        (target_type = 'post' AND post_id IS NOT NULL AND comment_id IS NULL) OR 
        (target_type = 'comment' AND comment_id IS NOT NULL AND post_id IS NULL)
    )
);
`
