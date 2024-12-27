-- categories
CREATE TABLE IF NOT EXISTS "category" (
    id         SERIAL PRIMARY KEY UNIQUE NOT NULL,
    slug       VARCHAR(255) UNIQUE       NOT NULL,
    title      VARCHAR(255) UNIQUE       NOT NULL,
    src        VARCHAR(255),
    alt        VARCHAR(255),
    parent_id  INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_parent_category FOREIGN KEY (parent_id) REFERENCES category (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "filter" (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    values      VARCHAR(255)[],        -- Array of strings
    filter_type VARCHAR(50)  NOT NULL, -- "equality", "range"
    min_value   FLOAT,                 -- For range filters
    max_value   FLOAT,                 -- For range filters
    category_id INTEGER      NOT NULL, -- fk for Category

    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
);

-- items
CREATE TABLE IF NOT EXISTS "item" (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    article     VARCHAR(255),
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

CREATE TABLE IF NOT EXISTS "item_label" (
    item_id UUID         NOT NULL,
    name    VARCHAR(255) NOT NULL, -- "hit", "new", "rec", etc

    PRIMARY KEY (item_id, name),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "item_category" (
    item_id     UUID    NOT NULL,
    category_id INTEGER NOT NULL,

    PRIMARY KEY (item_id, category_id),
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
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

-- cart
CREATE TABLE IF NOT EXISTS "cart" (
    id      SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS "cart_item" (
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quantity INTEGER NOT NULL,

    cart_id  BIGINT  NOT NULL,
    item_id  UUID    NOT NULL,
    CONSTRAINT fk_cart FOREIGN KEY (cart_id) REFERENCES cart (id) ON DELETE CASCADE,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
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
    price      DECIMAL(10, 2) NOT NULL,

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
    id          SERIAL PRIMARY KEY,
    slug        VARCHAR(255) UNIQUE NOT NULL,
    title       VARCHAR(255)        NOT NULL,
    description TEXT,
    src         VARCHAR(255)        NOT NULL,
    alt         VARCHAR(255),

    lasts_to    TIMESTAMP           NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "promotion_item" (
    discount     INTEGER NOT NULL,
    promotion_id INTEGER NOT NULL,
    item_id      UUID    NOT NULL,

    PRIMARY KEY (promotion_id, item_id),
    CONSTRAINT fk_promotion FOREIGN KEY (promotion_id) REFERENCES promotion (id) ON DELETE CASCADE,
    CONSTRAINT fk_item FOREIGN KEY (item_id) REFERENCES item (id) ON DELETE CASCADE
);