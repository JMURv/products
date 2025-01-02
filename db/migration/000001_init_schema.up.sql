-- categories
CREATE TABLE IF NOT EXISTS "category" (
    slug        VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    title       VARCHAR(255) UNIQUE             NOT NULL,
    src         VARCHAR(255),
    alt         VARCHAR(255),
    parent_slug VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_parent_category FOREIGN KEY (parent_slug) REFERENCES category (slug) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "filter" (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    values        VARCHAR(255)[],        -- Array of strings
    filter_type   VARCHAR(50)  NOT NULL, -- "equality", "range"
    min_value     FLOAT,                 -- For range filters
    max_value     FLOAT,                 -- For range filters
    category_slug VARCHAR(255) NOT NULL, -- fk for Category

    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_category FOREIGN KEY (category_slug) REFERENCES category (slug) ON DELETE CASCADE
);

-- items
CREATE TABLE IF NOT EXISTS "item" (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    article     VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    price       MONEY,
    src         VARCHAR(255),
    alt         VARCHAR(255),

    in_stock    BOOLEAN          DEFAULT TRUE,

    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,

    parent_id   UUID, -- self-referencing foreign key for variants
    CONSTRAINT fk_parent_item FOREIGN KEY (parent_id) REFERENCES item (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "label" (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
INSERT INTO "label" (name)
VALUES ('hit'),
       ('new'),
       ('rec');

CREATE TABLE IF NOT EXISTS "item_label" (
    item_id  UUID    NOT NULL,
    label_id INTEGER NOT NULL,

    PRIMARY KEY (item_id, label_id),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES "item" (id) ON DELETE CASCADE,
    CONSTRAINT fk_label FOREIGN KEY (label_id) REFERENCES "label" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "item_category" (
    item_id       UUID         NOT NULL,
    category_slug VARCHAR(255) NOT NULL,

    PRIMARY KEY (item_id, category_slug),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY (category_slug) REFERENCES category (slug) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "item_media" (
    id         SERIAL PRIMARY KEY,
    src        VARCHAR(255) NOT NULL,
    alt        VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    item_id    UUID         NOT NULL,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "item_attr" (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    value   TEXT         NOT NULL,

    item_id UUID         NOT NULL,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "related_product" (
    PRIMARY KEY (item_id, related_item_id),
    item_id         UUID NOT NULL,
    related_item_id UUID NOT NULL,

    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE,
    CONSTRAINT fk_related_item FOREIGN KEY (related_item_id) REFERENCES item (id) ON DELETE CASCADE
);

-- orders
CREATE TABLE IF NOT EXISTS "order" (
    id             SERIAL PRIMARY KEY,
    status         VARCHAR(50)    NOT NULL,
    total_amount   DECIMAL(10, 2) NOT NULL,
    fio            VARCHAR(255)   NOT NULL,
    tel            VARCHAR(20)    NOT NULL,
    email          VARCHAR(255)   NOT NULL,
    address        VARCHAR(255)   NOT NULL,
    delivery       VARCHAR(255)   NOT NULL,
    payment_method VARCHAR(255)   NOT NULL,
    user_id        UUID           NOT NULL,

    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "order_item" (
    id         SERIAL PRIMARY KEY,
    quantity   INTEGER        NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    order_id   BIGINT         NOT NULL,
    item_id    UUID           NOT NULL,
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES "order" (id) ON DELETE CASCADE,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);

-- favorites
CREATE TABLE IF NOT EXISTS "favorites" (
    user_id UUID NOT NULL,
    item_id UUID NOT NULL,

    PRIMARY KEY (user_id, item_id),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);

-- promotion
CREATE TABLE IF NOT EXISTS "promotion" (
    slug        VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    title       VARCHAR(255)                    NOT NULL,
    description TEXT,
    src         VARCHAR(255)                    NOT NULL,
    alt         VARCHAR(255),

    lasts_to    TIMESTAMP                       NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "promotion_item" (
    discount       INTEGER      NOT NULL,
    promotion_slug VARCHAR(255) NOT NULL,
    item_id        UUID         NOT NULL,

    PRIMARY KEY (promotion_slug, item_id),
    CONSTRAINT fk_promotion FOREIGN KEY (promotion_slug) REFERENCES promotion (slug) ON DELETE CASCADE,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);