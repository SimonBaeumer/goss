<?php
if (!isset($_SERVER['PHP_AUTH_USER'])) {
    header('WWW-Authenticate: Basic realm="My Realm"');
    header('HTTP/1.0 401 Unauthorized');
    echo 'not authorized!';
    exit;
} else {
    echo "user: {$_SERVER['PHP_AUTH_USER']}" . "</br>";
    echo "password: {$_SERVER['PHP_AUTH_PW']}";
}
