<?php

$headers = getallheaders();

if ($headers["Goss-Test"] === "worked!" && $headers["Another"] === "more") {
    echo "success";
    exit;
}

http_response_code(400);
echo "fail" . PHP_EOL;
echo "got headers: " . PHP_EOL;
