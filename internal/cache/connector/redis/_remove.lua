local function remove(object_id)
    local setKey = keySetKey(object_id)
    local keys = redis.call("SMEMBERS", setKey)
    local n = #keys
    for i = 1, n do
        redis.call("UNLINK", keys[i])
    end
    redis.call("UNLINK", setKey)
    redis.call("UNLINK", object_id)
end
