--你的验证码在 Redis 上的 Key
local key = KEYS[1]
--验证次数,一个验证码最多验证三次。这个记录还可以验证几次
local cntKey = key..":cnt"
--你的验证码 112233
local val = ARGV[1]
--过期时间
local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
    --key 存在，但是没有过期时间
    --系统错误，同事手残，手动设置了这个key ，但是没给过期时间
    return -2
    -- 540 = 600 -60 九分钟
elseif ttl == -2 or ttl < 540 then
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    -- 一切正常
    return 0
else
    return -1
end
