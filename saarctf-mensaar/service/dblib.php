<?php

require 'db-config.php';

class Db {

    /**
     * @var null|PDO
     */
    public static $db = null;

    public static function init() {
        if (!self::$db) {
            self::$db = new PDO(getConnectionConfig());
            self::$db->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        }
    }

    public static function query($query, $params = null) {
        if (!self::$db)
            self::init();
        if (!$params) {
            return self::$db->query($query)->fetchAll(PDO::FETCH_ASSOC);
        } else {
            $stmt = self::$db->prepare($query);
            foreach ($params as $k => $v) {
                $stmt->bindValue($k, $v);
            }
            $stmt->execute();
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        }
    }

    public static function queryList($query, $params = null) {
        if (!self::$db)
            self::init();
        if (!$params) {
            $data = self::$db->query($query)->fetchAll(PDO::FETCH_NUM);
        } else {
            $stmt = self::$db->prepare($query);
            foreach ($params as $k => $v) {
                $stmt->bindValue($k, $v);
            }
            $stmt->execute();
            $data = $stmt->fetchAll(PDO::FETCH_NUM);
        }
        $result = [];
        foreach ($data as $row) {
            if (count($row) == 1) {
                $result[] = $row[0];
            } else {
                $result[$row[0]] = $row[1];
            }
        }
        return $result;
    }

    public static function execute($query, $params = null) {
        if (!self::$db)
            self::init();
        if (!$params) {
            return self::$db->exec($query);
        } else {
            $stmt = self::$db->prepare($query);
			foreach ($params as $k => $v) {
                $stmt->bindValue($k, $v);
			}
			$result=$stmt->execute();
			return $result;   
		}
    }

    public function store_menu_db($menu) {
	$days = array("monday", "tuesday", "wednesday", "thursday", "friday");
    $menu = unserialize($menu);
	foreach ($days as $d) {
            store($menu, $d);
        }
    }

}
