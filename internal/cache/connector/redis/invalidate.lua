local n = #KEYS
for i = 1, n do
    local result = redis.call("GET", KEYS[i])
    -- a key that isn't cached (anymore) must not stop the invalidation of the remaining keys
    if result ~= false then
        remove(tostring(result))
    end
end
