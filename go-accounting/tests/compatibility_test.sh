#!/bin/bash

OLD_PORT=${1:-8003}
NEW_PORT=${2:-7002}
OLD_BASE="http://localhost:$OLD_PORT"
NEW_BASE="http://localhost:$NEW_PORT"

echo "=========================================="
echo "接口兼容性测试"
echo "Node.js 服务: $OLD_BASE"
echo "Go 服务: $NEW_BASE"
echo "=========================================="

compare_json() {
    local endpoint=$1
    local description=$2
    
    echo ""
    echo "--- $description ---"
    echo "Endpoint: $endpoint"
    
    old_response=$(curl -s "$OLD_BASE$endpoint" 2>/dev/null)
    new_response=$(curl -s "$NEW_BASE$endpoint" 2>/dev/null)
    
    echo "Node.js 响应:"
    echo "$old_response" | python3 -m json.tool 2>/dev/null || echo "$old_response"
    
    echo ""
    echo "Go 响应:"
    echo "$new_response" | python3 -m json.tool 2>/dev/null || echo "$new_response"
    
    old_keys=$(echo "$old_response" | python3 -c "import sys,json; d=json.load(sys.stdin); print(sorted([k for k in d.keys()]))" 2>/dev/null)
    new_keys=$(echo "$new_response" | python3 -c "import sys,json; d=json.load(sys.stdin); print(sorted([k for k in d.keys()]))" 2>/dev/null)
    
    if [ "$old_keys" == "$new_keys" ]; then
        echo "✅ 响应结构一致"
    else
        echo "❌ 响应结构不一致"
        echo "Node.js keys: $old_keys"
        echo "Go keys: $new_keys"
    fi
}

echo ""
echo "=========================================="
echo "1. 分类接口测试"
echo "=========================================="

compare_json "/api/categories" "获取所有分类"
compare_json "/api/categories?type=expense" "获取支出分类"

echo ""
echo "=========================================="
echo "2. 记账记录接口测试"
echo "=========================================="

compare_json "/api/records" "获取所有记录"
compare_json "/api/records/statistics" "获取统计信息"

echo ""
echo "=========================================="
echo "3. 创建分类测试"
echo "=========================================="

timestamp=$(date +%s)
test_category="{\"name\":\"测试分类_$timestamp\",\"type\":\"expense\"}"

echo "创建测试分类: $test_category"

old_create=$(curl -s -X POST -H "Content-Type: application/json" -d "$test_category" "$OLD_BASE/api/categories" 2>/dev/null)
new_create=$(curl -s -X POST -H "Content-Type: application/json" -d "$test_category" "$NEW_BASE/api/categories" 2>/dev/null)

echo "Node.js 创建响应:"
echo "$old_create" | python3 -m json.tool 2>/dev/null || echo "$old_create"

echo ""
echo "Go 创建响应:"
echo "$new_create" | python3 -m json.tool 2>/dev/null || echo "$new_create"

old_id=$(echo "$old_create" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('id',''))" 2>/dev/null)
new_id=$(echo "$new_create" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('id',''))" 2>/dev/null)

if [ -n "$old_id" ]; then
    curl -s -X DELETE "$OLD_BASE/api/categories/$old_id" > /dev/null 2>&1
fi
if [ -n "$new_id" ]; then
    curl -s -X DELETE "$NEW_BASE/api/categories/$new_id" > /dev/null 2>&1
fi

echo ""
echo "=========================================="
echo "4. 创建记录测试"
echo "=========================================="

test_record="{\"type\":\"expense\",\"category_id\":1,\"amount\":99.99,\"date\":\"2024-01-15\",\"note\":\"测试记录\"}"

echo "创建测试记录: $test_record"

old_record=$(curl -s -X POST -H "Content-Type: application/json" -d "$test_record" "$OLD_BASE/api/records" 2>/dev/null)
new_record=$(curl -s -X POST -H "Content-Type: application/json" -d "$test_record" "$NEW_BASE/api/records" 2>/dev/null)

echo "Node.js 创建记录响应:"
echo "$old_record" | python3 -m json.tool 2>/dev/null || echo "$old_record"

echo ""
echo "Go 创建记录响应:"
echo "$new_record" | python3 -m json.tool 2>/dev/null || echo "$new_record"

old_record_id=$(echo "$old_record" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('id',''))" 2>/dev/null)
new_record_id=$(echo "$new_record" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('id',''))" 2>/dev/null)

if [ -n "$old_record_id" ]; then
    curl -s -X DELETE "$OLD_BASE/api/records/$old_record_id" > /dev/null 2>&1
fi
if [ -n "$new_record_id" ]; then
    curl -s -X DELETE "$NEW_BASE/api/records/$new_record_id" > /dev/null 2>&1
fi

echo ""
echo "=========================================="
echo "测试完成"
echo "=========================================="
