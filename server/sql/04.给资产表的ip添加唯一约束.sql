-- 为方便进行资产管理，对同一个ip地址添加唯一约束
ALTER TABLE cs_asset
    ADD CONSTRAINT unique_asset_ip UNIQUE (asset_ip);