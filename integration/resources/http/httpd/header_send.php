<?php

const DO_NOT_REPLACE_HEADERS= false;

header("more: testing", DO_NOT_REPLACE_HEADERS);
header("more: duplicate", DO_NOT_REPLACE_HEADERS);

exit;