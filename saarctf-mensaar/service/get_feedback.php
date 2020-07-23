<?php

require_once "dblib.php";

set_include_path('lib/');

include('Crypt/RSA.php');

$db = new Db;
$db->init();
if ( isset($_GET['id']) ) {
    $feedback = $db->query("SELECT * FROM feedback WHERE cook_token=:token", array("token"=>$_GET['id']));
    if ( count($feedback) > 0 ) {
        $message = $feedback[0]['obj'];
        $pinned_key = "l\x85y\x17\xfd\xa6\xe4\xbbX\xf9^\xed~\xa1\xe9]\x16Q\xfd#s\x86M\n\xc0G\x87\xf0\xa9\x8e\x81q";
        $ciphertext = sodium_crypto_box_seal(
            $message,
            $pinned_key
        );
        echo $ciphertext;
    } else {
        echo "No feedback found :(";
    }
} else {
    echo "No token supplied :(";
}
