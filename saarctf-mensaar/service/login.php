<?php
require_once "dblib.php";

$db = new Db;
$db->init();

if (isset($_POST['mail']) && isset($_POST['pwd'])){
    $query_mail =array("email"=>$_POST['mail']);
    $result = $db->query("SELECT * from user_profile WHERE email = :email LIMIT 1", $query_mail);
    if (!empty($result) && isset($result[0]['pwd']) && password_verify($_POST['pwd'], $result[0]['pwd'])){
        session_start();
        if (isset($result[0]['email']))
            $_SESSION['email'] = $result[0]['email'];
        if (isset($result[0]['name']))
            $_SESSION['name'] = $result[0]['name'];
        if (isset($_POST['remember'])){
            $token = hash( 'sha3-256' , $result[0]['pwd'].$_SESSION['name']);
            $json = array("email"=>$_SESSION['email'], "token"=>$token);
            setcookie("member_me", base64_encode(json_encode($json)), time() +(10 * 365 * 24 * 60 * 60));
        }
    } else {
        http_response_code(500);
    }
    header("Location: index.php");
    die();
}
