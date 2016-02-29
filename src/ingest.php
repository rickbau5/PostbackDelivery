<?php

function braced($str) {
    return '{' . $str . '}';
}

$postdata = file_get_contents('php://input');
$decoded = json_decode($postdata,true);

if ($decoded) {
    if (isset($decoded['data']) && isset($decoded['endpoint'])) {
        // Verify this
        $endpointURL = $decoded['endpoint']['url'];
        $endpointMethod = $decoded['endpoint']['method'];

        foreach($decoded['data'] as &$data) {
            $key = 'key';
            $value = 'value';

            if (isset($data[$key]) && isset($data[$value])) {    
                $keyed = str_replace(braced($key), $data[$key], $endpointURL);
                $populated = str_replace(braced($value), $data[$value], $keyed);

                var_dump($data);
                echo $populated . PHP_EOL;
            } else {
                echo 'Malformed data object.' . PHP_EOL;
            }
        }
    } else {
        echo 'No data received.' . PHP_EOL;
    }
} else {
   echo "Sorry, no.";
}
?>
