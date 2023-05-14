CREATE TABLE IF NOT EXISTS articles (
                                     id int(11) unsigned NOT NULL AUTO_INCREMENT,
                                     title varchar(255) NOT NULL,
                                     body Text NOT NULL,
                                     created_at datetime DEFAULT CURRENT_TIMESTAMP,
                                     updated_at datetime DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (id)
) ENGINE=InnoDB;