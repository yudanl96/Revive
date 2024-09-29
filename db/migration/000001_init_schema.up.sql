CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()), 
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE  NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE posts (
  id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()), 
  user_id VARCHAR(36) NOT NULL,
  description VARCHAR(255) NOT NULL,
  price INT DEFAULT 0 NOT NULL,
  sold boolean DEFAULT 0 NOT NULL,  -- Use 0 (false) instead of 'False'
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- Tracks updates
  FOREIGN KEY (user_id) REFERENCES users(id) 
);

CREATE TABLE saved_posts (
    user_id VARCHAR(36) NOT NULL,                        -- The user who saved the post
    post_id VARCHAR(36) NOT NULL,                        -- The post that was saved
    saved_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp of when the post was saved
    PRIMARY KEY (user_id, post_id),             -- Composite primary key to ensure unique saves
    FOREIGN KEY (user_id) REFERENCES users(id), -- Relationship between saved posts and users
    FOREIGN KEY (post_id) REFERENCES posts(id)  -- Relationship between saved posts and posts
);
