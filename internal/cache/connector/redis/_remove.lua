local function remove(object_id)
    local setKey = keySetKey(object_id)
    local keys = redis.call("SMEMBERS", setKey)
    local n = #keys
    for i = 1, n do
        redis.call("DEL", keys[i])
    end
    redis.call("DEL", setKey)
    redis.call("DEL", object_id)
end
