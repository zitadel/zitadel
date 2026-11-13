-- keySetKey returns the redis key of the set containing all keys to the object.
local function keySetKey (object_id)
    return object_id .. "-keys"
end

-- resetTTL sets the TTL of the object, its key set and all keys still pointing to the object.
local function resetTTL(object_id, ttl)
    redis.call("EXPIRE", object_id, ttl)
    local setKey = keySetKey(object_id)
    local keys = redis.call("SMEMBERS", setKey)
    local n = #keys
    for i = 1, n do
        -- a key overwritten by a newer object is managed by that object
        if redis.call("GET", keys[i]) == object_id then
            redis.call("EXPIRE", keys[i], ttl)
        end
    end
    redis.call("EXPIRE", setKey, ttl)
end

local function getTime()
    return tonumber(redis.call('TIME')[1])
end

-- getCall wrapts redis.call so a nil is returned instead of false.
local function getCall (...)
    local result = redis.call(...)
    if result == false then
        return nil
    end
    return result 
end
