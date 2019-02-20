#!/usr/bin/env bats

export GOSS_PATH=./../../../release/goss-linux-amd64
export GOSS_FILES_STRATEGY=mount

@test "run dgoss with mount strategy" {
    run ./../dgoss run nginx:latest
    echo "$output"
    [ "$status" -eq 0 ]
    [ "${lines[0]}" == "INFO: Starting docker container" ]
}

@test "run dgoss with cp strategy" {
    export GOSS_FILES_STRATEGY=cp
    run ./../dgoss run nginx:latest
    echo "$output"
    [ "$status" -eq 0 ]
    [ "${lines[0]}" == "INFO: Creating docker container" ]
    [ "${lines[1]}" == "INFO: Copy goss files into container" ]
}

@test "run dgoss with invalid goss path" {
    GOSS="/invalid"
    run GOSS_PATH=${GOSS} ./../dgoss run nginx:latest
    [ "$status" -eq 127 ]
}