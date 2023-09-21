ALTER TABLE projections.quotas ALTER COLUMN from_anchor DROP NOT NULL;
ALTER TABLE projections.quotas ALTER COLUMN amount DROP NOT NULL;
ALTER TABLE projections.quotas ALTER COLUMN interval DROP NOT NULL;
ALTER TABLE projections.quotas ALTER COLUMN limit_usage DROP NOT NULL;
