# admin
> pkusz  
> pkusz

# blockchain
> mis  
> mis20201001

# MongoDB设置密码
```shell script
mongo
use admin
db.createUser({user: 'pkusz', pwd: 'pkusz', roles: ['root']})
exit
sudo vim /etc/mongod.conf
# change
    #security:
# to
    security:
        authorization: enabled
sudo service mongod restart
mongo
use admin
db.auth("pkusz","pkusz")
use blockchain
db.createUser({ user: "mis", pwd: "mis20201001", roles: [{ role: "dbOwner", db: "blockchain" }] })
```

# MongoDB设置索引
在区块链启动之后，会在本地mongo数据库中建立以-LOG结尾的日志collection
如果日志中的数据量过大会导致取数据排序过程中内存溢出
因此需要设置排序的字段为索引，只需要执行下面一行命令

```shell script
db.getCollection("xxxxxxxxxx-LOG").createIndex({"timestamp":1})
```