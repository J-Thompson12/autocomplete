!#/bin/sh

set -e -u

go test ./... -coverprofile cover.out -covermode atomic

perc=`go tool cover -func=cover.out | tail -n 1 | sed -Ee 's!^[^[:digit:]]+([[:digit:]]+(\.[[:digit:]]+))%$!\1!'`
res=`echo "$perc >= 80.0" | bc`
test "$res" -eq 1 && exit 0
echo "Insufficient coverage: $perc" >&2
exit 1