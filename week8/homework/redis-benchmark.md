# pprof

[TOC]

## 名词解释

### info memory

1. used_memory：由 Redis 分配器分配的内存总量，包含了redis进程内部的开销和数据占用的内存，以字节（byte）为单位

2. used_memory_rss：向操作系统申请的内存大小

3. used_memory_peak：redis的内存消耗峰值(以字节为单位)

4. used_memory_peak_perc：使用内存达到峰值内存的百分比，即(used_memory / used_memory_peak) * 100%

5. used_memory_overhead：Redis为了维护数据集的内部机制所需的内存开销，包括所有客户端输出缓冲区、查询缓冲区、AOF重写缓冲区和主从复制的backlog

6. used_memory_startup：Redis服务器启动时消耗的内存

7. used_memory_dataset：数据占用的内存大小，即used_memory - used_memory_overhead

8. used_memory_dataset_perc：数据占用的内存大小的百分比，

   100% * (used_memory_dataset / (used_memory - used_memory_startup))

9. total_system_memory：整个系统内存

10. used_memory_lua：Lua脚本存储占用的内存

11. maxmemory：Redis实例的最大内存配置

12. maxmemory_policy：当达到maxmemory时的淘汰策略

13. mem_fragmentation_ratio：碎片率，used_memory_rss / used_memory

14. mem_allocator：内存分配器

15. active_defrag_running：表示没有活动的defrag任务正在运行，1表示有活动的defrag任务正在运行（defrag:表示内存碎片整理）

16. lazyfree_pending_objects：0表示不存在延迟释放的挂起对象

### redis-benchmark

1. -h: 指定服务器主机名（默认127.0.0.1）
2. -p: 指定服务器端口（默认6379）
3. -s: 指定服务器 socket
4. -c: 指定并发连接数（默认50）
5. -n: 指定请求数（默认10000）
6. -d: 以字节的形式指定 SET/GET 值的数据大小（默认2）
7. -k: 1=keep alive, 0=reconnect（默认1）
8. -r: SET/GET/INCR 使用随机 key, SADD 使用随机值
9. -p: 通过管道传输 <numreq> 请求 （默认1）
10. -q: 强制退出 redis。仅显示 query/sec 值
11. --csv: 以 CSV 格式输出
12. -l: 生成循环，永久执行测试
13. -t: 仅运行以逗号分隔的测试命令列表
14. -i: Idle 模式。仅打开 N 个 idle 连接并等待

## value = 10

```shell
# Memory before
used_memory:842608
used_memory_human:822.86K
used_memory_rss:11710464
used_memory_rss_human:11.17M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:10474
used_memory_dataset_perc:17.43%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:13.90
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 10 -t get,set -q
SET: 138580.94 requests per second
GET: 126968.01 requests per second
```

```shell
# Memory after
used_memory:18103408
used_memory_human:17.26M
used_memory_rss:24154112
used_memory_rss_human:23.04M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.38%
used_memory_overhead:9247486
used_memory_startup:782504
used_memory_dataset:8855922
used_memory_dataset_perc:51.13%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.33
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

# DBSIZE
(integer) 157955 # ? 压测25w随机key，为啥只有16w的key呢？随机数重复导致的么？
```

``` python
# ((used_memory_dataset_after - used_memory_dataset_before) - (dbsize * d)) / dbsize
>>> rs = 8855922 - 10474
>>> rs
8845448
>>> rs -= 157955 * 10
>>> rs
7265898
>>> rs /= 157955
>>> rs
45.999797410654935
```

## value = 20

```shell
# Memory before
used_memory:842784
used_memory_human:823.03K
used_memory_rss:6307840
used_memory_rss_human:6.02M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:10650
used_memory_dataset_perc:17.67%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:7.48
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 20 -t get,set -q
SET: 135648.39 requests per second
GET: 126198.89 requests per second
```

```shell
# Memory after
used_memory:20638144
used_memory_human:19.68M
used_memory_rss:26132480
used_memory_rss_human:24.92M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.44%
used_memory_overhead:9250086
used_memory_startup:782504
used_memory_dataset:11388058
used_memory_dataset_perc:57.35%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.27
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

# DBSIZE
(integer) 158020
```

```python
>>> rs = 11388058 - 10650
>>> rs
11377408
>>> rs -= 158020 * 20
>>> rs
8217008
>>> rs /= 158020
>>> rs
51.999797493988105
```

## value = 50

```shell
# Memory before
used_memory:842960
used_memory_human:823.20K
used_memory_rss:6975488
used_memory_rss_human:6.65M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:10826
used_memory_dataset_perc:17.91%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:8.27
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 50 -t get,set -q
SET: 135062.12 requests per second
GET: 127811.86 requests per second
```

```shell
# Memory after
used_memory:25691360
used_memory_human:24.50M
used_memory_rss:31252480
used_memory_rss_human:29.80M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.54%
used_memory_overhead:9249086
used_memory_startup:782504
used_memory_dataset:16442274
used_memory_dataset_perc:66.01%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.22
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

# DBSIZE
(integer) 157995
```

```python
>>> rs = 16442274 - 10826
>>> rs
16431448
>>> rs -= 157995 * 50
>>> rs
8531698
>>> rs /= 157995
>>> rs
53.999797461945
```

## value = 100

```shell
# Memory before
used_memory:843136
used_memory_human:823.38K
used_memory_rss:7266304
used_memory_rss_human:6.93M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:11002
used_memory_dataset_perc:18.15%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:8.62
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 100 -t get,set -q
SET: 134192.16 requests per second
GET: 129265.77 requests per second
```

```shell
# Memory after
used_memory:33206752
used_memory_human:31.67M
used_memory_rss:38768640
used_memory_rss_human:36.97M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.70%
used_memory_overhead:9234806
used_memory_startup:782504
used_memory_dataset:23971946
used_memory_dataset_perc:73.93%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.17
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

#DBSIZE
(integer) 157638
```

```python
>>> rs = 23971946 - 11002
>>> rs
23960944
>>> rs -= 157638 * 100
>>> rs
8197144
>>> rs /= 157638
>>> rs
51.99979700326063
```

## value = 200

```shell
# Memory before
used_memory:843312
used_memory_human:823.55K
used_memory_rss:7888896
used_memory_rss_human:7.52M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:11178
used_memory_dataset_perc:18.38%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:9.35
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 200 -t get,set -q
SET: 134770.89 requests per second
GET: 126262.62 requests per second
```

```shell
# Memory after
used_memory:51005872
used_memory_human:48.64M
used_memory_rss:57139200
used_memory_rss_human:54.49M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:1.08%
used_memory_overhead:9253686
used_memory_startup:782504
used_memory_dataset:41752186
used_memory_dataset_perc:83.13%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.12
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

# DBSIZE
(integer) 158110
```

```python
>>> rs = 41752186 - 11178
>>> rs
41741008
>>> rs -= 158110 * 200
>>> rs
10119008
>>> rs /= 158110
>>> rs
63.999797609259375
```

## value = 1000

```shell
# Memory before
used_memory:843488
used_memory_human:823.72K
used_memory_rss:8224768
used_memory_rss_human:7.84M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:11354
used_memory_dataset_perc:18.62%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:9.75
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 1000 -t get,set -q
SET: 130890.05 requests per second
GET: 131164.73 requests per second
```

```shell
# Memory after
used_memory:177474176
used_memory_human:169.25M
used_memory_rss:185737216
used_memory_rss_human:177.13M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:3.75%
used_memory_overhead:9252966
used_memory_startup:782504
used_memory_dataset:168221210
used_memory_dataset_perc:95.21%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.05
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

# DBSIZE
(integer) 158092
```

```python
>>> rs = 168221210 - 11354
>>> rs
168209856
>>> rs -= 158092 * 1000
>>> rs
10117856
>>> rs /= 158092
>>> rs
63.99979758621562
```

## value = 5000

```shell
# Memory before
used_memory:843840
used_memory_human:824.06K
used_memory_rss:9740288
used_memory_rss_human:9.29M
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:0.02%
used_memory_overhead:832134
used_memory_startup:782504
used_memory_dataset:11706
used_memory_dataset_perc:19.09%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:11.54
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0
```

```shell
# redis-benchmark -r 250000 -n 250000 -c 1000 -d 5000 -t get,set -q
SET: 106564.37 requests per second
GET: 113327.28 requests per second
```

```shell
# Memory after
used_memory:1307319552
used_memory_human:1.22G
used_memory_rss:1336999936
used_memory_rss_human:1.25G
used_memory_peak:4728498944
used_memory_peak_human:4.40G
used_memory_peak_perc:27.65%
used_memory_overhead:9236726
used_memory_startup:782504
used_memory_dataset:1298082826
used_memory_dataset_perc:99.35%
total_system_memory:33599131648
total_system_memory_human:31.29G
used_memory_lua:37888
used_memory_lua_human:37.00K
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
mem_fragmentation_ratio:1.02
mem_allocator:jemalloc-3.6.0
active_defrag_running:0
lazyfree_pending_objects:0

# DBSIZE
(integer) 157686
```

```python
>>> rs = 1298082826 - 11706
>>> rs
1298071120
>>> rs -= 157686 * 5000
>>> rs
509641120
>>> rs /= 157686
>>> rs
3231.9997970650534
```

