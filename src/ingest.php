<?php
$postdata = file_get_contents('php://input');
$decoded = json_decode($postdata,true);

if ($decoded) {
    var_dump($decoded);

    echo "\n";

    if (isset($decoded['data'])) {
        foreach($decoded['data'] as &$data) {
            var_dump($data);
        }
    }
} else {
   echo "Sorry, no.";
}
?>
