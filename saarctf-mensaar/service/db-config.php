<?php

require_once "objects.php";
require_once "dblib.php";
$dbConnectionConfig = "pgsql:dbname=mensaar;user=mensaar";

function getConnectionConfig() {
    global $dbConnectionConfig;
    return $dbConnectionConfig;
}

function json_parse($x) {
    return unserialize($x);
};

function store_menu($food,$day)
{
	$db = new Db;
	$db->init();
        $data = array("food"=>$food, "day"=>$day);
        $db->execute("INSERT INTO menu (food, day, date) VALUES (:food, :day, now())", $data);
};

function store($menu, $day) {
	$fnc = $menu->func['save'];
	foreach ($menu->$day as $f) {
		$fnc($f, $day);
	}
}
?>
