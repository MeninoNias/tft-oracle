ALTER TABLE items ADD COLUMN type TEXT NOT NULL DEFAULT 'item';
CREATE INDEX idx_items_type ON items(type);
