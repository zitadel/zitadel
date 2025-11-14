-- SELECT ensures the DB namespace for each script.
-- When used, it consumes the first ARGV entry.
redis.call("SELECT", ARGV[1])
