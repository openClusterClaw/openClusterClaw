#!/bin/bash

# Happy 重启脚本

echo "=== Happy 重启脚本 ==="
echo ""

# 查找并停止 Happy 进程
echo "1. 停止 Happy 进程..."
HAPPY_PIDS=$(pgrep -f "happy-coder" | xargs)
if [ -n "$HAPPY_PIDS" ]; then
    echo "找到进程: $HAPPY_PIDS"
    pkill -f "happy-coder"
    sleep 2
    echo "✓ Happy 进程已停止"
else
    echo "✓ 没有运行中的 Happy 进程"
fi

# 等待进程完全停止
sleep 1

# 启动 Happy
echo ""
echo "2. 启动 Happy..."
happy start &
HAPPY_PID=$!

# 等待启动
sleep 2

# 检查是否成功启动
if ps -p $HAPPY_PID > /dev/null 2>&1; then
    echo "✓ Happy 已启动 (PID: $HAPPY_PID)"
    echo ""
    echo "Happy 进程状态:"
    ps aux | grep -i happy | grep -v grep
else
    echo "✗ Happy 启动失败"
    exit 1
fi

echo ""
echo "=== Happy 重启完成 ==="