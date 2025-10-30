# 
# RedisShakeCMD 快速使用
- 命令帮助
```
$ ./redis-shake-cmd aof_reader -h
------redis-shake-cmd aof_reader (description)------
Usage:
  aof_reader [flags]

Flags:
  -f, --filepath string   [required] (default "/tmp/.aof")
  -h, --help              help for aof_reader
  -a, --timestamp int     # subsecond
```
- 从aof中筛选出指定前缀的命令,并导入到redis
```
$ ./redis-shake-cmd aof_reader -f /tmp/appendonly.aof filter -p TEST: file_writer -f /tmp/cmd.txt
$ redis-cli -p 6379 < /tmp/cmd.txt
```
- 从一个redis离线导出aof再导入到另一个redis
```
$ ./redis-shake-cmd scan_reader -c -d 0 -n 1024 -a 172.33.66.164:8300 -u default -p 123456 file_writer -t aof -f /tmp/a.aof
$ redis-cli -p 6379 -a 123456 --pipe < /tmp/a.aof
```
- 把rdb/aof文件导入到运行的redis中
```
$ ./redis-shake-cmd rdb_reader -f /tmp/dump.rdb redis_writer -a 172.33.66.164:6379 -u default -p 123456 advanced --ncpu=1
$ ./redis-shake-cmd aof_reader -f /tmp/appendonly.aof redis_writer -a 172.33.66.164:6379 -u default -p 123456 advanced --ncpu=1
```
- 扫描一个redis存储到另一个redis
```
$ ./redis-shake-cmd scan_reader -c -a 172.33.66.164:8300 -u default -p 123456  redis_writer -c -a 172.33.66.117:8300
```
# 完整帮助命令
```
$ ./bin/redis-shake-cmd -h
redis-shake-cmd [command_reader][flags] [command_writer][flags] [filter][flags] [advanced][flags] [module]
command_reader: aof_reader, rdb_reader, scan_reader, sync_reader, sync_reader.sentinel
command_writer: redis_writer, redis_writer.sentinel


------redis-shake-cmd aof_reader (description)------

Usage:
  aof_reader [flags]

Flags:
  -f, --filepath string   [required] (default "/tmp/.aof")
  -h, --help              help for aof_reader
  -a, --timestamp int     # subsecond

------redis-shake-cmd rdb_reader (description)------

Usage:
  rdb_reader [flags]

Flags:
  -f, --filepath string   [required] (default "/tmp/dump.rdb")
  -h, --help              help for rdb_reader

------redis-shake-cmd scan_reader (description)------

Usage:
  scan_reader [flags]

Flags:
  -a, --address string    # [required] when cluster is true, set address to one of the cluster node (default "127.0.0.1:6379")
  -c, --cluster           # set to true if source is a redis cluster
  -n, --count int         # number of keys to scan per iteration (default 1)
  -d, --dbs ints          # set you want to scan dbs such as [1,5,7], if you don't want to scan all
  -h, --help              help for scan_reader
  -k, --ksn               # set to true to enabled Redis keyspace notifications (KSN) subscription
  -p, --password string   # keep empty if no authentication is required
  -r, --prefer_replica    
  -s, --scan              # set to false if you don't want to scan keys (default true)
  -t, --tls               
  -u, --username string   # keep empty if not using ACL

------redis-shake-cmd sync_reader (description)------

Usage:
  sync_reader [flags]

Flags:
  -a, --address string    # [required]For clusters, specify the address of any cluster node; use the master or slave address in master-slave mode (default "127.0.0.1:6379")
  -c, --cluster           # Set to true if the source is a Redis cluster
  -h, --help              help for sync_reader
  -p, --password string   # Keep empty if no authentication is required
  -r, --prefer_replica    # Set to true to sync from a replica node
  -o, --sync_aof          # Set to false if AOF synchronization is not required (default true)
  -d, --sync_rdb          # Set to false if RDB synchronization is not required (default true)
  -t, --tls               # Set to true to enable TLS if needed
  -l, --try_diskless      # Set to true for diskless sync if the source has repl-diskless-sync=yes
  -u, --username string   # Keep empty if ACL is not in use

------redis-shake-cmd sync_reader.sentinel (description)------

Usage:
  sync_reader.sentinel [flags]

Flags:
  -a, --address string       [required]eg: 127.0.0.1:6379
  -h, --help                 help for sync_reader.sentinel
  -m, --master_name string   
  -p, --password string      
  -t, --tls                  
  -u, --username string

------redis-shake-cmd redis_writer (description)------

Usage:
  redis_writer [flags]

Flags:
  -a, --address string    # when cluster is true, set address to one of the cluster node (default "127.0.0.1:6380")
  -c, --cluster           # set to true if target is a redis cluster
  -h, --help              help for redis_writer
  -o, --off_reply         # turn off the server reply
  -p, --password string   # keep empty if no authentication is required
  -t, --tls               # turn off the server reply
  -u, --username string   # keep empty if not using ACL

------redis-shake-cmd redis_writer.sentinel (description)------

Usage:
  redis_writer.sentinel [flags]

Flags:
  -a, --address string       [required]eg: 127.0.0.1:6379
  -h, --help                 help for redis_writer.sentinel
  -m, --master_name string   
  -p, --password string      
  -t, --tls                  
  -u, --username string

------redis-shake-cmd filter (description)------

Usage:
  filter [flags]

Flags:
  -m, --allow_command strings         # Allow or block specific commands
                                      # Examples:
                                      #   allow_command = ["GET", "SET"]  # Only allow GET and SET commands
                                      #   block_command = ["DEL", "FLUSHDB"]  # Block DEL and FLUSHDB commands
                                      # Leave empty to allow all commands
  -g, --allow_command_group strings   # Allow or block specific command groups
                                      # Available groups:
                                      #   SERVER, STRING, CLUSTER, CONNECTION, BITMAP, LIST, SORTED_SET,
                                      #   GENERIC, TRANSACTIONS, SCRIPTING, TAIRHASH, TAIRSTRING, TAIRZSET,
                                      #   GEO, HASH, HYPERLOGLOG, PUBSUB, SET, SENTINEL, STREAM
                                      # Examples:
                                      #   allow_command_group = ["STRING", "HASH"]  # Only allow STRING and HASH commands
                                      #   block_command_group = ["SCRIPTING", "PUBSUB"]  # Block SCRIPTING and PUBSUB commands
                                      # Leave empty to allow all command groups
  -w, --allow_db ints                 # Specify allowed and blocked database numbers (e.g., allow_db = [0, 1, 2], block_db = [3, 4, 5])
                                      # Leave empty to allow all databases
  -p, --allow_key_prefix strings      
  -e, --allow_key_regex strings       
  -s, --allow_key_suffix strings      
  -a, --allow_keys strings            # Allow keys with specific prefixes or suffixes
                                      # Examples:
                                      #   allow_keys = ["user:1001", "product:2001"]
                                      #   allow_key_prefix = ["user:", "product:"]
                                      #   allow_key_suffix = [":active", ":valid"]
                                      #   allow A collection of keys containing 11-digit mobile phone numbers
                                      #   allow_key_regex = [":\\d{11}:"]
                                      # Leave empty to allow all keys
  -l, --block_command strings         
  -u, --block_command_group strings   
  -o, --block_db ints                 
  -r, --block_key_prefix strings      
  -x, --block_key_regex strings       
  -i, --block_key_suffix strings      
  -k, --block_keys strings            # Block keys with specific prefixes or suffixes
                                      # Examples:
                                      #   block_keys = ["temp:1001", "cache:2001"]
                                      #   block_key_prefix = ["temp:", "cache:"]
                                      #   block_key_suffix = [":tmp", ":old"]
                                      #   block test 11-digit mobile phone numbers keys
                                      #   block_key_regex = [":test:\\d{11}:"]
                                      # Leave empty to block nothing
  -f, --function string               # Function for custom data processing
                                      # For best practices and examples, visit:
                                      # https://tair-opensource.github.io/RedisShake/zh/filter/function.html
  -h, --help                          help for filter

------redis-shake-cmd advanced (description)------

Usage:
  advanced [flags]

Flags:
  -w, --aws_psync string                            # If the source is Elasticache, you can set this item. AWS ElastiCache has custom
                                                    # psync command, which can be obtained through a ticket.
  -d, --dir string                                   (default "data")
  -y, --empty_db_before_sync                        # destination will delete itself entire database before fetching files
                                                    # from source during full synchronization.
                                                    # This option is similar redis replicas RDB diskless load option:
                                                    #   repl-diskless-load on-empty-db
  -h, --help                                        help for advanced
  -l, --log_file string                              (default "shake.log")
  -i, --log_interval int                            # in seconds (default 5)
  -e, --log_level string                            # debug, info or warn (default "info")
  -a, --ncpu int                                    # runtime.GOMAXPROCS, 0 means use runtime.NumCPU() cpu cores
  -f, --pprof_port int                              # pprof port, 0 means disable
  -b, --rdb_restore_command_behavior string         # panic, rewrite or skip
                                                    # redis-shake-cmd gets key and value from rdb file, and uses RESTORE command to
                                                    # create the key in target redis. Redis RESTORE will return a "Target key name
                                                    # is busy" error when key already exists. You can use this configuration item
                                                    # to change the default behavior of restore:
                                                    # panic:   redis-shake-cmd will stop when meet "Target key name is busy" error.
                                                    # rewrite: redis-shake-cmd will replace the key with new value.
                                                    # skip:  redis-shake-cmd will skip restore the key when meet "Target key name is busy" error. (default "panic")
  -t, --status_port int                             # status port, 0 means disable
  -q, --target_redis_client_max_querybuf_len uint   # This setting corresponds to the 'client-query-buffer-limit' in Redis configuration.
                                                    # The default value is typically 1GB.
                                                    # It's recommended not to modify this value unless absolutely necessary. (default 1073741824)
  -x, --target_redis_proto_max_bulk_len int         # This setting corresponds to the 'proto-max-bulk-len' in Redis configuration.
                                                    # It defines the maximum size of a single string element in the Redis protocol.
                                                    # The value must be 1MB or greater. Default is 512MB.
                                                    # It's recommended not to modify this value unless absolutely necessary. (default 512000000)

------redis-shake-cmd module (description)------

Usage:
  module [flags]

Flags:
  -h, --help                         help for module
  -v, --target_mbbloom_version int   # The data format for BF.LOADCHUNK is not compatible in different versions. v2.6.3 <=> 20603

------redis-shake-cmd file_writer (description)------

Usage:
  file_writer [flags]

Flags:
  -f, --filepath string   [required] (default "/tmp/cmd.txt")
  -h, --help              help for file_writer
  -t, --type string       # default: cmd, options: cmd/json/aof (default "cmd")

```
# 源码改造来自
[RedisShake](https://github.com/tair-opensource/RedisShake/)
