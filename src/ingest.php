<?php

function braced($str) {
    return '{' . $str . '}';
}

$postdata = file_get_contents('php://input');
$decoded = json_decode($postdata, true);

if ($decoded) {
    $redis = new Redis();
    $connected = $redis->connect('127.0.0.1');

    if ($connected) {
        if (isset($decoded['data']) && isset($decoded['endpoint'])) {
            $redis->lPush('requests', $postdata);
        } else {
            echo 'No data received.' . PHP_EOL;
        }
    } else {
        echo 'Unable to connect to DB' . PHP_EOL;
    }
} else {
   echo "Sorry, no.";
}
?>
