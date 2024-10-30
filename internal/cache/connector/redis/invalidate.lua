local n = #KEYS
for i = 1, n do
    local result = redis.call("GET", KEYS[i])
    if result == false then
        return nil
    end
    local object_id = tostring(result)
    remove(object_id)
end
