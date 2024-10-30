-- keySetKey returns the redis key of the set containing all keys to the object.
local function keySetKey (object_id)
    return object_id .. "-keys"
end

local function getTime()
    return tonumber(redis.call('TIME')[1])
end
