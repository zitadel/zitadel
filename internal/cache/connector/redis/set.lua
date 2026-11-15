-- KEYS: [1]: object_id; [>1]: index keys.
local object_id = KEYS[1]
local object = ARGV[2]
local usage_lifetime = tonumber(ARGV[3]) -- usage based lifetime in seconds
local max_age = tonumber(ARGV[4]) -- max age liftime in seconds

redis.call("HSET", object_id,"object", object)
if usage_lifetime > 0 then
    redis.call("HSET", object_id, "usage_lifetime", usage_lifetime)
    -- enable usage based TTL
    redis.call("EXPIRE", object_id, usage_lifetime)
    if max_age > 0 then
        -- set max_age to hash map for expired remove on Get
        local expiry = getTime() + max_age
        redis.call("HSET", object_id, "expiry", expiry)
    end
elseif max_age > 0 then
    -- enable max_age based TTL
    redis.call("EXPIRE", object_id, max_age)
end

-- Index keys and the key set need a TTL as well. Otherwise they are never removed when the object
-- isn't requested again by that key, for example after a version upgrade,
-- and can't be evicted under a volatile-* maxmemory-policy.
-- The object is never returned after max_age, so the keys can simply expire then.
-- Without max_age they share the usage based TTL of the object, which Get resets.
local keys_ttl = max_age
if keys_ttl <= 0 then
    keys_ttl = usage_lifetime
end

local n = #KEYS
local setKey = keySetKey(object_id)
for i = 2, n do -- offset to the second element to skip object_id
    redis.call("SADD", setKey, KEYS[i]) -- set of all keys used for housekeeping
    if keys_ttl > 0 then
        redis.call("SET", KEYS[i], object_id, "EX", keys_ttl) -- key to object_id mapping
    else
        redis.call("SET", KEYS[i], object_id) -- key to object_id mapping
    end
end
if keys_ttl > 0 then
    redis.call("EXPIRE", setKey, keys_ttl)
end
