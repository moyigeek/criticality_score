#!/bin/bash

# 设定工作目录和输出CSV文件
WORK_DIR="./extracted"
OUTPUT_FILE="packages.csv"

# 创建或清空输出CSV文件
echo "Package Name,URL" > "$OUTPUT_FILE"

# 遍历所有包名文件夹
for package_dir in "$WORK_DIR"/*; do
    if [ -d "$package_dir" ]; then
        package_name=$(basename "$package_dir")
        desc_file="$package_dir/desc"
        
        # 检查desc文件是否存在
        if [ -f "$desc_file" ]; then
            # 获取URL行
            url=$(awk '/^%URL%/{getline; print}' "$desc_file")
            
            # 检查url是否非空
            if [ -n "$url" ]; then
                # 将包名和URL写入CSV文件
                echo "$package_name,$url" >> "$OUTPUT_FILE"
            fi
        fi
    fi
done

echo "CSV文件已生成: $OUTPUT_FILE"

