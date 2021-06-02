ALTER TABLE adminapi.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE adminapi.features ADD COLUMN label_policy_watermark BOOLEAN;
ALTER TABLE auth.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE auth.features ADD COLUMN label_policy_watermark BOOLEAN;
ALTER TABLE authz.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE authz.features ADD COLUMN label_policy_watermark BOOLEAN;
ALTER TABLE management.features RENAME COLUMN label_policy TO label_policy_private_label;
ALTER TABLE management.features ADD COLUMN label_policy_watermark BOOLEAN;