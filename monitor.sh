#!/bin/bash

# 设置日志文件路径
LOG_FILE="/tmp/memory_usage.log"

# 获取当前内存使用情况 (以KB为单位)
current_usage=$(ps -o rss --noheaders -p $(pgrep app))

# 获取当前时间戳
timestamp=$(date +%s)

# 记录当前内存使用情况
echo "$timestamp $current_usage" >>$LOG_FILE

# 保留最近的 10 个记录
tail -n 10 $LOG_FILE >$LOG_FILE.tmp && mv $LOG_FILE.tmp $LOG_FILE

# 计算内存使用增长率
readings=$(wc -l <$LOG_FILE)
if [ "$readings" -ge 2 ]; then
	start_time=$(head -n 1 $LOG_FILE | awk '{print $1}')
	start_usage=$(head -n 1 $LOG_FILE | awk '{print $2}')
	end_time=$(tail -n 1 $LOG_FILE | awk '{print $1}')
	end_usage=$(tail -n 1 $LOG_FILE | awk '{print $2}')

	time_diff=$((end_time - start_time))
	usage_diff=$((end_usage - start_usage))

	# 增长率 (KB/秒)
	growth_rate=$(echo "$usage_diff / $time_diff" | bc -l)
else
	growth_rate=0
fi

# 设置阈值 (KB/秒)
threshold=10

# 检查增长率是否超过阈值
if [ "$(echo "$growth_rate > $threshold" | bc -l)" -eq 1 ]; then
	echo "CRITICAL - Memory growth rate is $growth_rate KB/sec"
	exit 2
else
	echo "OK - Memory growth rate is $growth_rate KB/sec"
	exit 0
fi
