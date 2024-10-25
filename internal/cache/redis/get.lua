local result = redis.call("GET", KEYS[1])
if result == false then
    return nil
end
local object_id = tostring(result)

local entries = redis.call("HGETALL", object_id)
if entries == nil then
    -- object expired, but there are keys that need to be cleaned up
    remove(object_id)
    return nil
end

-- entries is a key-value paired string array.
-- Extract the values into variables we understand.
local object = entries[2]
local usage_lifetime = tonumber(entries[4])
local expiry = tonumber(entries[6])

-- max-age must be checked manually
if not (expiry == nil) and expiry > 0 then
    if getTime() > expiry then
        remove(object_id)
        return nil
    end
end
-- reset usage based TTL
if not (usage_lifetime == nil) and usage_lifetime > 0 then
    redis.call('EXPIRE', object_id, usage_lifetime)
end

return object
