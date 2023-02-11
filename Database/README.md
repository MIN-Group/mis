# admin
> pkusz  
> pkusz

# blockchain
> mis  
> mis20201001

# Set password of MongoDB
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

# Set index of MongoDB
After the blockchain is started, a log collection ending with -LOG will be established in the local mongo database.
If the amount of data in the log is too large, it will cause memory overflow during the data fetching and sorting process.
Therefore, it is necessary to set the sorted field as an index, and only need to execute the following line of command.

```shell script
db.getCollection("xxxxxxxxxx-LOG").createIndex({"timestamp":1})
```