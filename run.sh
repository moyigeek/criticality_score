AWK_SCRIPT=$(cat <<'EOF'
BEGIN {
    cmd = "go run cmd/home2git/main.go"
    appendix = ""  # > to file or other commands
}
/^$/ {
    pkg = ""
}
/^Package: / {
    pkg = gensub("^Package: ", "", "g", $0)
}
/^Homepage: / {
    if (pkg != "") {
        printf "%s -package %s -homepage %s %s\n", cmd, pkg, gensub("^Homepage: ", "", "g", $0), appendix
    }
}
EOF
)
curl https://mirrors.hust.edu.cn/debian/dists/stable/main/source/Sources.gz 2>/dev/null |
	gunzip - |
	awk "$AWK_SCRIPT" | bash -x

