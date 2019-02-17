#!/usr/bin/env bash

#!/usr/bin/env bash

function add_group() {
    groupadd test-group
}

function add_user() {
    useradd test-user
}

function add_file() {
    echo "Create file..."
    echo "test-content" > /tmp/test-file
    chmod 0644 /tmp/test-file

    echo "Create symlink"
    ln -s /tmp/test-file symlink
}

add_file
add_group
add_user