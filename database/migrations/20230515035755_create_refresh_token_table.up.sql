CREATE TABLE IF NOT EXISTS refresh_tokens (
                                     id int(11) unsigned NOT NULL AUTO_INCREMENT,
                                     token TEXT NOT NULL,
                                     user_id int NOT NULL,
                                     created_at datetime DEFAULT CURRENT_TIMESTAMP,
                                     updated_at datetime DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (id)
) ENGINE=InnoDB;