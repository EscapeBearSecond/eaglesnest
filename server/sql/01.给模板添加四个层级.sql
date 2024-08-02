ALTER TABLE cs_template
ADD COLUMN tag1 text NOT NULL DEFAULT '',
ADD COLUMN tag2 text NOT NULL DEFAULT '',
ADD COLUMN tag3 text NOT NULL DEFAULT '',
ADD COLUMN tag4 text NOT NULL DEFAULT '';

COMMENT ON COLUMN cs_template.tag1 IS '资产大类';
COMMENT ON COLUMN cs_template.tag2 IS '资产小类';
COMMENT ON COLUMN cs_template.tag3 IS '厂商名称';
COMMENT ON COLUMN cs_template.tag4 IS '型号';