-- KEYS: [1]: object_id; [>1]: index keys.
local object_id = KEYS[1]
local object = ARGV[2]
local usage_lifetime = tonumber(ARGV[3]) -- usage based lifetime in seconds
local max_age = tonumber(ARGV[4]) -- max age liftime in seconds

-- The object, its key set and the index keys share the same TTL.
-- Index keys without TTL would never be removed when the object isn't requested again by that key,
-- for example after a version upgrade, and can't be evicted under a volatile-* maxmemory-policy.
local ttl = 0
if usage_lifetime > 0 then
    ttl = usage_lifetime -- usage based TTL, reset on Get
elseif max_age > 0 then
    ttl = max_age -- max_age based TTL
end

redis.call("HSET", object_id,"object", object)
if ttl > 0 then
    redis.call("EXPIRE", object_id, ttl)
end
if usage_lifetime > 0 then
    redis.call("HSET", object_id, "usage_lifetime", usage_lifetime)
    if max_age > 0 then
        -- set max_age to hash map for expired remove on Get
        local expiry = getTime() + max_age
        redis.call("HSET", object_id, "expiry", expiry)
    end
end

local n = #KEYS
local setKey = keySetKey(object_id)
for i = 2, n do -- offset to the second element to skip object_id
    redis.call("SADD", setKey, KEYS[i]) -- set of all keys used for housekeeping
    if ttl > 0 then
        redis.call("EXPIRE", setKey, ttl)
        redis.call("SET", KEYS[i], object_id, "EX", ttl) -- key to object_id mapping
    else
        redis.call("SET", KEYS[i], object_id) -- key to object_id mapping
    end
end
