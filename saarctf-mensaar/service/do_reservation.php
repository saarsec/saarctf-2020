<?php

require_once "dblib.php";

$db = new Db;
$db->init();

session_start();
if (isset($_POST['reset'])) {
    $db->execute("UPDATE seat SET last_reservation=NULL, reserved_by='' WHERE last_reservation < NOW() - INTERVAL '6 minutes';");
    http_response_code(200);
    die();
}
if (isset($_SESSION['name']) && isset($_SESSION['email']) && session_status() === PHP_SESSION_ACTIVE) {

    $usermail = $_SESSION['email'];
    $res_seat = $db->query("SELECT * FROM seat WHERE reserved_by=:mail", array("mail" => $usermail));

    if (count($res_seat) === 0) {
        if (isset($_POST['row']) && isset($_POST['table']) && isset($_POST['seatrow']) && isset($_POST['seatnumber']) && isset($_POST['user']) && $usermail === $_POST['user']) {
            if ($_POST['row'] === '0' && $_POST['table'] === '0' && $_POST['seatrow'] === '0') {
                $picnics = $db->query("SELECT * FROM seat WHERE row=0 AND tablenumber=0 AND seatrow=0 AND reserved_by='';");
                if (0 < count($picnics)) {
                    $db->execute("UPDATE seat SET reserved_by=:mail, last_reservation=NOW() WHERE row=0 AND tablenumber=0 AND seatrow=0 AND seatnumber=:seatnumber", array("mail" => $_POST['user'], "seatnumber" => $picnics[0]['seatnumber']));
                    http_response_code(200);
                    die();
                } else {
                    $max_picnic = $db->query("SELECT MAX(seatnumber) FROM seat WHERE row=0 AND tablenumber=0 AND seatrow=0");
                    $db->execute("INSERT INTO seat VALUES (0, 0, 0, :seatnumber, :mail, NOW())", array("mail" => $_POST['user'], "seatnumber" => $max_picnic[0][0] + 1));
                    http_response_code(200);
                    die();
                }
            } else {
                $seat_status = $db->query("SELECT * FROM seat WHERE row=:row AND tablenumber=:tablenumber AND seatrow=:seatrow AND seatnumber=:seatnumber AND reserved_by=''", array("row" => $_POST['row'], "tablenumber" => $_POST['table'], "seatrow" => $_POST['seatrow'], "seatnumber" => $_POST['seatnumber']));
                if (0 < count($seat_status)) {
                    $db->execute("UPDATE seat SET reserved_by=:mail, last_reservation=NOW() WHERE row=:row AND tablenumber=:tablenumber AND seatrow=:seatrow AND seatnumber=:seatnumber", array("mail" => $_POST['user'], "row" => $_POST['row'], "tablenumber" => $_POST['table'], "seatrow" => $_POST['seatrow'], "seatnumber" => $_POST['seatnumber']));
                    http_response_code(200);
                    die();
                } else {
                    echo "Seat is not free";
                    http_response_code(500);
                    die();
                }
            }
        } else {
            echo "Request destroyed";
            http_response_code(500);
            die();
        }
    } else {
        echo "User already has a seat";
        http_response_code(500);
        die();
    }
}

session_destroy();
http_response_code(500);
die();


