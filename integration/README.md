# Integration tests

## Set `GOSS_EXE` env

```bash
export GOSS_EXE=/../../../release/goss-linux-amd64
```

## Commander

The tool `commander` is used to run the test suites.
You find more about it [here](https://github.com/SimonBaeumer/commander).

```bash
# run tests
$ commander test

# execute a single test
$ commander test commander.yaml "kernel-param resouce should fail"

# verbose execution for debugging
$ commander test --verbose 
```

## Debugging shell scripts

Add `-x` and `-v` to the shebang lines, this enables printing of the executed commands.

```bash
#!/bin/bash -vx
```