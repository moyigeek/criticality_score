AWK_SCRIPT=$(cat <<'EOF'
BEGIN {
    OFS = ","  # 设置输出字段分隔符为逗号，用于CSV格式
}
/^$/ {
    pkg = ""  # Reset package name at the beginning of a new record
}
/^Package: / {
    pkg = gensub("^Package: ", "", "g", $0)  # Extract package name
}
/^Homepage: / {
    if (pkg != "") {
        homepage = gensub("^Homepage: ", "", "g", $0)  # Extract homepage URL
        if (!(homepage in seen)) {
            seen[homepage] = 1
            if (homepage ~ /(https?|git):\/\/(www\.)?github\.com\/[^\/]+\/?$/) {
                printf "%s,%s,%s\n", pkg, homepage, homepage  # Output package, homepage, and detected Git link
            } else if (homepage ~ /(https?|git):\/\/[^\/]*git[^\/]*\/[^\/]+\/[^\/]+/) {
                printf "%s,%s,%s\n", pkg, homepage, homepage
            } else {
                printf "%s,%s,\n", pkg, homepage  # No Git link detected
            }
        }
    }
}
EOF
)

# 现在，使用上述脚本执行 curl 和 awk 命令。
curl https://mirrors.hust.edu.cn/debian/dists/stable/main/source/Sources.gz 2>/dev/null |
    gunzip - |
    awk "$AWK_SCRIPT" > results.csv

