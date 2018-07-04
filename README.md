[![Go Report Card](https://goreportcard.com/badge/github.com/barryz/goci)](https://goreportcard.com/report/github.com/barryz/goci)
![GoCI](http://goci.ele.me/na/goci/eleme/goci/badge?type=job)
[![Build Status](https://travis-ci.org/barryz/goci.svg?branch=master)](https://travis-ci.org/barryz/goci)

**goci** is a command-line tool for checking the code quality of Go locally.

It supports:

- `build`  project build
- `gofmt`  gofmt checking
- `govet`  code quality checking
- `golint` code style checking
- `test`   go test or unit test
- `race`   race condition test


## Installation
```
go get -u github.com/barryz/goci
```

## Run
```
goci -c goci.yml
```

## Configuration
Create a file which named `goci.yml`. This file should include fields as below:

| fields   |  type  | comment |
|----------|--------|-------|
| build    | string | Command or script that used for build project. Scripts should use the relative path. eg: ./ |
| test     | string | Command for testing |
| excludes | array  | The directories which in excludes will not be fmt, lint or ver |
| race     | struct | Execute race condition testing |
| lint     | struct | Configurations for golint |
| skips    | array  | steps which to skipped |



##### `race` struct fields

| fields |  type | comment |
|---------|--------|-------|
| main    | string | entry-point file for project，eg：main.go |
| main_cmd | string | command for execution |
| script  | string | test script eg: sh race_test.sh |
| timeout | int    | timeout for race condition execution |

##### `lint` struct fields

| fields |  type  | comment  |
|---------|--------|-------|
| ignore_no_comment_error | bool | false(default)|


Example:

```yaml
build: go build
test: go test
excludes:
    - templates # except templates
    - vendor # except vendor
race:
    # go run -race main.go
    main: ./main.go
    # kill the above after 20 seconds
    timeout: 20
    # run this script to interact with the above running process (as a test)
    script: ./ab.sh
    # after 20 seconds, if no race detected, the goci/race passs
skips:
    - fmt
```