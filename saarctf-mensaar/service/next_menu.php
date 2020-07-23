<?php

require_once "objects.php";
require_once "dblib.php";

$db = new Db;
$db->init();

if ( isset($_POST['data']) ) {
	$message = $_POST['token'];
	$key_bytes= array(0xee,0xad,0xef,0x10,0xe1,0xa3,0x2d,0xce,0x3f,0x20,0x52,0x5b,0x1f,0xc6,0x8d,0x1a,0x35,0x1c,0x6f,0xd4,0xd3,0x41,0xd5,0xb0,0xb6,0x10,0x9a,0x7a,0xb4,0x4f,0xd0,0x76);
	foreach ($key_bytes as $key => $value){
		$key_bytes[$key] = chr($value);
	}
	$pinned_key = implode("", $key_bytes);
	$massage = sodium_crypto_sign_open(
		$message,
		$pinned_key
	);
	if (sha1($massage) === $_POST['hash']) {
		$db->store_menu_db($_POST['data']);
	} else {
		echo "Invalid signature, "; 
    }
} else {
    error_log("Missing Data for next menu.");
}
