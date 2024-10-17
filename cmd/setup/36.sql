SELECT instance_id, type, reached_date, last_pushed_date, primary_domain
FROM projections.milestones
WHERE reached_date IS NOT NULL
ORDER BY instance_id, reached_date;
