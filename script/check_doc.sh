#!/usr/bin/bash

ls -1 docs/rules | sort -n > expected
go run rules/rule_names/main.go > actual

diff expected actual