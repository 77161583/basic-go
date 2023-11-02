local key = KEYS[1]
--用户输入的code
local expectedCode = ARGV[1]
local code = redis.call("get",key)
local cntKey = key..":cnt"
--转成一个数字
local cnt = tonumber(redis.call("get",cntKey))
if cnt == nil or cnt <= 0 then
    --说明用户一直输错,有人在搞破坏
    --或者已经用过了，也是有人在搞你
    return -1
elseif expectedCode == code then
    --输入对了
    --用完不能在用
    redis.call("set",cntKey,-1)
    return 0
else
    --用户输入错误了
    --可验证次数 -1
    redis.call("decr",cntKey,-1)
    return -2
end