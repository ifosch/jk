#!/bin/env -S bash -l

set -eu

tests_check() {
    echo "*** Run unit tests"
    go test -v ./... -race -covermode=atomic
}

complexity_check() {
    echo "*** Complexity check"
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    gocyclo -avg -top 5 -over 15 .
}

format_check() {
    echo "*** Format check"
    gofmt -s -e -d -l . | tee /tmp/gofmt.output && [ $(cat /tmp/gofmt.output | wc -l) -eq 0 ]
}

inefficiencies_check() {
    echo "*** Inefficiencies check"
    go install github.com/gordonklaus/ineffassign@latest
    go mod tidy
    ineffassign ./...
}

smells_check() {
    echo "*** Smells check"
    go mod tidy
    go vet ./...
}

spelling_check() {
    echo "*** Spelling check"
    go install github.com/client9/misspell/cmd/misspell@latest
    misspell -error .
}

static_check() {
    echo "*** Static check"
    go install honnef.co/go/tools/cmd/staticcheck@latest
    go mod download
    staticcheck ./...
}

style_check() {
    echo "*** Style check"
    go install golang.org/x/lint/golint@latest
    golint ./...
}

check_functions() {
    declare -F | awk '{print $3}' | grep -E "_check$" | sed -e 's/_check$//'
}

try() {
    if check_functions | grep -w ${1} >/dev/null; then
	${1}_check && echo "=== OK!" || (echo "=== NOK!" && return -1)
    else
        echo No ${1}_check available
	return 255
    fi
}

if [ "${1}" == "all" ]; then
    failure=0
    set +e
    for check in $(check_functions); do
        try ${check} || failure=1
    done
    set -e
    exit ${failure}
else
    try ${1}
fi
