-- keySetKey returns the redis key of the set containing all keys to the object.
local function keySetKey (object_id)
    return object_id .. "-keys"
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
