# mis-blockchain

[中文版见此（Chinese version see here）](#jump)


- #### Introduction

  mis-blockchain is a blockchain platform based on the PPoV (Parallel Proof of Vote) algorithm, developed in GO language, and supports identity management, node management, logging and other functions of the multi-identifier network system MIN

- #### Run

  -  Quick start

  ```shell script
  go build -o MIS-BC
  ./MIS-BC
  ```
  >This method is to start 4 blockchain nodes locally

   -  General form

  ```shell script
  ./MIS-BC -f <config_file>
  ```

  >You can modify the configuration file according to the actual situation or use the configuration file generation tool under common/config.go to generate a new configuration file. If it fails to run, please check whether the MongoDB database service is enabled. If it is not enabled, you can use the following command to enable it:                                                                                                                                      

  ```shell script
  sudo service mongod start
  ```

  >If MongoDB has not been installed, please refer to the following webpage to install it:                                                                                                                                       

  <https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/>


- #### Installation of development environment
  - go mod environment
  ```shell script
   1.Set GOMODULE on
     		go env -w GO111MODULE="on"
   2.Update code
     After pulling the latest code, Goland will prompt you to detect Go Moudule. You need to click Enabled. Note: The following environment can be filled in or not. If necessary, fill in GOPROXY="https://proxy.golang.com.cn"
   3.Change Go Proxy
     		go env -w GOPROXY="https://proxy.golang.com.cn"
     Then type in
     		go env
     To check status
   4.update go mod
     		go mod tidy
  ```

  -  GO dependencies

  ```shell script
    # go mod will automatically detect mod files and install dependencies
    go get github.com/Hyperledger-TWGC/ccs-gm
    go get github.com/JodeZer/mgop 
    go get github.com/google/uuid 
    go get github.com/karlseguin/ccache/v2 
    go get github.com/larspensjo/config 
    go get github.com/tinylib/msgp 
    go get github.com/yudeguang/ratelimit 
    go get gopkg.in/alexcesaro/quotedprintable.v3 
    go get gopkg.in/check.v1 
    go get gopkg.in/gomail.v2 
    go get gopkg.in/mgo.v2 
    go get gopkg.in/yaml.v2
    ...... 
  ```

- #### Instructions for use

  By default (quick start), four blockchain nodes are running locally, with an IP address of 127.0.0.1, blockchain communication ports `5010, 5011, 5012, 5013`, and service ports `8010, 8011, 8012, 8013` for the management front-end.
  The connection establishment ports provided for vpn-management-server are `9999, 10000, 10001, 10002`, and the ports for normal communication are `6666, 6667, 6668, 6669`.

- #### Module division

  ```textmate
    AccountManager		|	Account management module
       Message			|	Message management module
       MetaData			|	Definition of metadata format
       Database			|	Database module
       Network			|	Network module
      Blockchain			|	Program core module
    TransactionPool		|	Transaction pool module
       security			|	Security module
       common				|	Configuration file management module
       utils		    	|	Toolkit module
  ```

- #### Configuration file description

  ```textmate
  [Log]
  "LogToFile"                 Whether to output the log to a file   
  "LogPath"                   Log file path
  "Level"                     Log display level: Panic 0,Fatal 1,Error 2,Warn 3,Info 4,Debug 5,Trace 6
  
  [Node]
  "WorkerList"                IP and port of bookkeepers
  "WorkerCandidateList"       IP and port of candidate bookkeepers
  "VoterList"                 IP and port of voters
  "BcManagementServerList"    Servers of blockchain management
  "ServerNum"	              	Servers number of blockchain management
  "SingleServerNodeNum"       Number of nodes running on this machine
  "IP"                        Local IP address
  "PubIP"                     Local public network address (Set to empty if none)
  "Port"                      Local port
  "Hostname"		          	Local node name
  "Areaname"		          	Name of the region where the node is located
  "Countryname"		        Name of the country where the node is located
  "Longitude"                 Longitude of the node's location
  "Latitude"                  Latitude of the node's location
  "CacheTime"                 Duration of caching node status information, in minutes
  "IsNewJoin"                 Indicates whether it is a newly joined node
  
  [Consensus] 
  "PubkeyList"                Public key of the node running on this machine
  "PrikeyList"                Private key of the node running on this machine
  "MyPubkey"                  Reserved configuration
  "MyPrikey"                  Reserved configuration
  "GenesisDutyWorker"         Node number that generated the genesis block
  "WorkerNum"                 Number of bookkeepers
  "VotedNum"                  Number of voters
  "BlockGroupPerCycle"        Rotation cycle of bookkeepers
  "Tcut"                      Overtime time
  "GenerateBlockPeriod"       Length of Generate block cycle
  "TxPoolSize"                Size of transaction pool
  
  [MIR]
  "SqlitePath"                File path sqlite database
  
  [SESSION]
  "DefaultExpiration"         Default validity period of the session, in minutes
  "CleanupInterval"           Cleaning cycle of a session, in minutes
  ```



*<span id="jump">中文版</span>*

- #### 介绍

  mis-blockchain是基于PPoV（Parallel Proof of Vote）算法的区块链平台，以GO语言开发而成，支持多标识网络体系MIN的身份管理、节点管理、日志记录等功能

- #### 运行

  -  快速启动

  ```shell script
  go build -o MIS-BC
  ./MIS-BC
  ```
  >该方式是在本地启动4个区块链节点

   -  一般形式

  ```shell script
  ./MIS-BC -f <config_file>
  ```

  >可按实际情况修改配置文件或使用common/config.go下的配置文件生成工具生成新的配置文件。如果无法运行请先检查MongoDB数据库服务是否开启，如未开启可使用以下命令开启:                                                                                                                                       

  ```shell script
  sudo service mongod start
  ```

  >如果MongoDB尚未安装，请参考如下网页进行安装:                                                                                                                                       

  <https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/>


- #### 开发环境安装
  - go mod环境
  ```shell script
   1.开启GOMODULE
     go env -w GO111MODULE="on"
   2.更新代码
     拉取最新代码后，Goland会提示你检测到Go Moudule,点击Enabled即可。 注意：之后的environment可填可不填,需要填的话填写GOPROXY="https://goproxy.io"
   3.修改Go Proxy
     go env -w GOPROXY="https://goproxy.io"
     再输入
     go env
     查看状态
   4.更新go mod
     go mod tidy
  ```

  -  使用到的GO插件

  ```shell script
    # go mod会自动检测mod文件并安装依赖
    go get github.com/Hyperledger-TWGC/ccs-gm
    go get github.com/JodeZer/mgop 
    go get github.com/google/uuid 
    go get github.com/karlseguin/ccache/v2 
    go get github.com/larspensjo/config 
    go get github.com/tinylib/msgp 
    go get github.com/yudeguang/ratelimit 
    go get gopkg.in/alexcesaro/quotedprintable.v3 
    go get gopkg.in/check.v1 
    go get gopkg.in/gomail.v2 
    go get gopkg.in/mgo.v2 
    go get gopkg.in/yaml.v2
    ...... 
  ```

- #### 使用说明

  默认（快速启动）情况下，在本地运行四个区块链节点，IP地址为127.0.0.1，区块链通信端口`5010, 5011, 5012, 5013`,对管理前端提供服务端口`8010, 8011, 8012, 8013`，
对vpn-management-server提供的建立连接端口为`9999，10000，10001，10002`，同时正常通信的端口为`6666，6667，6668，6669`

- #### 模块划分

  ```textmate
    AccountManager		|	账号管理模块
       Message			|	消息管理模块
       MetaData			|	元数据格式定义
       Database			|	数据库模块
       Network			|	网络模块
      Blockchain		|	程序核心模块
    TransactionPool		|	事务池模块
       security			|	安全模块
       common			|	配置文件管理模块
       utils		    |	工具包模块
  ```

- #### 配置文件说明

  ```textmate
  [Log]
  "LogToFile"                 是否输出日志到文件   
  "LogPath"                   日志文件路径
  "Level"                     日志显示等级Panic 0,Fatal 1,Error 2,Warn 3,Info 4,Debug 5,Trace 6
  
  [Node]
  "WorkerList"                记账节点IP与端口
  "WorkerCandidateList"       候选记账节点IP与端口
  "VoterList"                 投票节点IP与端口
  "BcManagementServerList"    区块链管理后台服务器列表
  "ServerNum"	              区块链管理后台服务器数量
  "SingleServerNodeNum"       本机运行的节点数
  "IP"                        本机IP地址
  "PubIP"                     本机公网地址（没有可空）
  "Port"                      本机端口
  "Hostname"		          本机节点名
  "Areaname"		          节点所在地区名
  "Countryname"		          节点所在国家名
  "Longitude"                 节点所在位置经度
  "Latitude"                  节点所在位置纬度
  "CacheTime"                 缓存节点状态信息的时长，单位为分钟
  "IsNewJoin"                 指示是否为新加入节点
  
  [Consensus] 
  "PubkeyList"                本机运行节点的公钥
  "PrikeyList"                本机运行节点的私钥
  "MyPubkey"                  保留字段
  "MyPrikey"                  保留字段
  "GenesisDutyWorker"         生成创世区块的节点编号
  "WorkerNum"                 记账节点数
  "VotedNum"                  投票节点数
  "BlockGroupPerCycle"        记账节点轮换周期
  "Tcut"                      超时时间
  "GenerateBlockPeriod"       产生区块周期
  "TxPoolSize"                交易池大小
  
  [MIR]
  "SqlitePath"                sqlite数据库文件地址
  
  [SESSION]
  "DefaultExpiration"         Session默认有效期，单位为分钟
  "CleanupInterval"           Session清理周期，单位为分钟
  ```