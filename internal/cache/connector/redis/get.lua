local result = redis.call("GET", KEYS[1])
if result == false then
    return nil
end
local object_id = tostring(result)

local object = getCall("HGET", object_id, "object")
if object == nil then
    -- object expired, but there are keys that need to be cleaned up
    remove(object_id)
    return nil
end

-- max-age must be checked manually
local expiry = getCall("HGET", object_id, "expiry")
if not (expiry == nil) and expiry > 0 then
    if getTime() > expiry then
        remove(object_id)
        return nil
    end
end

local usage_lifetime = getCall("HGET", object_id, "usage_lifetime")
-- reset usage based TTL
if not (usage_lifetime == nil) and tonumber(usage_lifetime) > 0 then
    redis.call('EXPIRE', object_id, usage_lifetime)
end

return object
