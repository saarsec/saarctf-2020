<?php
session_start();
session_destroy();
setcookie("member_me", "", time()-3600);
header("Location: index.php");
die();