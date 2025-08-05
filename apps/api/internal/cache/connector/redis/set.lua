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

local n = #KEYS
local setKey = keySetKey(object_id)
for i = 2, n do -- offset to the second element to skip object_id
    redis.call("SADD", setKey, KEYS[i]) -- set of all keys used for housekeeping
    redis.call("SET", KEYS[i], object_id) -- key to object_id mapping
end
