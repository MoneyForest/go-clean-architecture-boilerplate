CREATE TABLE IF NOT EXISTS matching (
    id CHAR(36) NOT NULL,
    me_id CHAR(36) NOT NULL,
    partner_id CHAR(36) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_matching_me_id (me_id),
    INDEX idx_matching_partner_id (partner_id),
    FOREIGN KEY (me_id) REFERENCES user(id),
    FOREIGN KEY (partner_id) REFERENCES user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
