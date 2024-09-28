CREATE TABLE `accounts` (
  `id` varchar(36) PRIMARY KEY,  -- UUID should be handled during INSERT
  `username` varchar(255) NOT NULL,
  `saved_id` varchar(36),
  `post_id` varchar(36),
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Use CURRENT_TIMESTAMP for now()
);

CREATE TABLE `posts` (
  `id` varchar(36) PRIMARY KEY,  -- UUID should be handled during INSERT
  `user_id` varchar(36) NOT NULL,
  `description` varchar(255) NOT NULL,
  `sold` boolean DEFAULT 0 NOT NULL,  -- Use 0 (false) instead of 'False'
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- Tracks updates
);

ALTER TABLE `accounts` 
  ADD FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`);
