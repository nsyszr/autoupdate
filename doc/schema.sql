CREATE SCHEMA tenant1;

-- update target
CREATE SEQUENCE tenant1.update_targets_id_seq;
CREATE TABLE tenant1.update_targets
(
    id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('tenant1.update_targets_id_seq'),
    title TEXT NOT NULL UNIQUE,
    slug VARCHAR(256) NOT NULL UNIQUE,
    download_list_id BIGINT,
    is_enabled BOOLEAN
);
ALTER SEQUENCE tenant1.update_targets_id_seq
OWNED BY tenant1.update_targets.id;

-- download list
CREATE SEQUENCE tenant1.download_lists_id_seq;
CREATE TABLE tenant1.download_lists
(
    id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('tenant1.download_lists_id_seq'),
    title TEXT NOT NULL UNIQUE,
    default_update_packet_id BIGINT
);
ALTER SEQUENCE tenant1.download_lists_id_seq
OWNED BY tenant1.download_lists.id;

-- download list line item
CREATE SEQUENCE tenant1.download_list_line_items_id_seq;
CREATE TABLE tenant1.download_list_lines_item
(
    id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('tenant1.download_list_line_items_id_seq'),
    download_list_id BIGINT NOT NULL,
    serial_number VARCHAR(200) NOT NULL UNIQUE,
    update_packet_id BIGINT NOT NULL
);
ALTER SEQUENCE tenant1.download_list_line_items_id_seq
OWNED BY tenant1.download_list_line_items.id;
