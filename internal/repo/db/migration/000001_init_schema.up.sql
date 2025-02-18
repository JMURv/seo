-- pages
CREATE TABLE IF NOT EXISTS page (
    slug       VARCHAR(255) PRIMARY KEY,
    title      VARCHAR(255),
    href       VARCHAR(255),

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- SEO
CREATE TABLE IF NOT EXISTS seo (
    id             SERIAL PRIMARY KEY,
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    keywords       TEXT,
    og_title       VARCHAR(255),
    og_description TEXT,
    og_image       VARCHAR(255),
    obj_name       VARCHAR(255) NOT NULL,
    obj_pk         VARCHAR(255) NOT NULL,

    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_seo_name ON seo (obj_name);
CREATE INDEX idx_seo_pk ON seo (obj_pk);
