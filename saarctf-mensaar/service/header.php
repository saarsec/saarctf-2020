<!DOCTYPE html>
<html lang="en">

<head>

    <?php

    require_once "dblib.php";

    $db = new Db;
    $db->init();

    session_start();
    if (!(isset($_SESSION['name']) && isset($_SESSION['email']) && session_status() === PHP_SESSION_ACTIVE) && array_key_exists("member_me", $_COOKIE)) {
        $json = json_decode(base64_decode($_COOKIE["member_me"]),true);
        $user_token = $json["token"];
        $mail = $json["email"];
        if (isset( $user_token ) && isset( $mail )) {
            $user_data =  $db->query("SELECT * FROM user_profile WHERE email=:email", array("email"=>$mail));
            if (count($user_data) > 0) {
                $real_token = hash( 'sha3-256' , $user_data[0]['pwd'].$user_data[0]['name']);
                if ($user_token == $real_token) {
                        $_SESSION['email'] = $user_data[0]['email'];
                        $_SESSION['name'] = $user_data[0]['name'];
                } else {
                    session_destroy();
                }
            } else {
                session_destroy();
            }
        } else {
            session_destroy();
        }
    }

    ?>

    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>Mensa of the Saarland</title>

    <!-- Bootstrap core CSS -->
    <link href="vendor/bootstrap/css/bootstrap.css" rel="stylesheet">

    <!-- Bootstrap core JavaScript -->
    <script src="vendor/jquery/jquery.min.js"></script>
    <script src="vendor/bootstrap/js/bootstrap.bundle.min.js"></script>

    <style>

        #login-dp{
            min-width: 250px;
            padding: 14px 14px 0;
            overflow:hidden;
            background-color:rgba(255,255,255,.8);
        }
        #login-dp .help-block{
            font-size:12px
        }
        #login-dp .bottom{
            width:100%;
            background-color:rgba(255,255,255,.8);
            border-top:1px solid #ddd;
            clear:both;
            padding:14px;
        }
        #login-dp .social-buttons{
            margin:12px 0
        }
        #login-dp .social-buttons a{
            width: 49%;
        }
        #login-dp .form-group {
            margin-bottom: 10px;
        }
        .btn-fb{
            color: #fff;
            background-color:#3b5998;
        }
        .btn-fb:hover{
            color: #fff;
            background-color:#496ebc
        }
        .btn-tw{
            color: #fff;
            background-color:#55acee;
        }
        .btn-tw:hover{
            color: #fff;
            background-color:#59b5fa;
        }
        @media(max-width:768px){
            #login-dp{
                background-color: inherit;
                color: #fff;
            }
            #login-dp .bottom{
                background-color: inherit;
                border-top:0 none;
            }
        }

        .rotate90 {
            -webkit-transform: rotate(90deg);
            -moz-transform: rotate(90deg);
            -o-transform: rotate(90deg);
            -ms-transform: rotate(90deg);
            transform: rotate(90deg);
        }

        .rotate270 {
            -webkit-transform: rotate(270deg);
            -moz-transform: rotate(270deg);
            -o-transform: rotate(270deg);
            -ms-transform: rotate(270deg);
            transform: rotate(270deg);
        }


    </style>
</head>

<body>

<!-- Navigation -->
<nav class="navbar navbar-expand-lg navbar-dark bg-dark fixed-top">
    <div class="container">
        <a class="navbar-brand" href="index.php">Mensa of the Saarland</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarResponsive" aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarResponsive">
            <ul class="navbar-nav ml-auto">
                <li class="nav-item<?php if ($page === 'home') echo " active"; ?>">
                    <a class="nav-link" href="index.php">Home
                        <?php if ($page === 'home') echo "<span class='sr-only'>(current)</span>"; ?>
                    </a>
                </li>
                <?php
                if (isset($_SESSION['name']) && isset($_SESSION['email']) && session_status() === PHP_SESSION_ACTIVE) {
                    ?>
                    <li class="nav-item<?php if ($page === 'reservation') echo " active"; ?>">
                        <a class="nav-link" href="reservation.php">Reservation
                            <?php if ($page === 'reservation') echo "<span class='sr-only'>(current)</span>"; ?>
                        </a>
                    </li>
                    <li class="nav-item<?php if ($page === 'feedback') echo " active"; ?>">
                        <a class="nav-link" href="feedback.php">Feedback
                            <?php if ($page === 'feedback') echo "<span class='sr-only'>(current)</span>"; ?>
                        </a>
                    </li>
                    <li class="nav-item<?php if ($page === 'foodDB') echo " active"; ?>">
                        <a class="nav-link" href="foodDB.php">Food Database
                            <?php if ($page === 'foodDB') echo "<span class='sr-only'>(current)</span>"; ?>
                        </a>
                    </li>
                    <?php
                }
                ?>
                <li class="nav-item<?php if ($page === 'about') echo " active"; ?>">
                    <a class="nav-link" href="about.php">About
                        <?php if ($page === 'about') echo "<span class='sr-only'>(current)</span>"; ?>
                    </a>
                </li>
            </ul>
            <div style="width:25%;"></div>
            <?php
            if (isset($_SESSION['name']) && isset($_SESSION['email']) && session_status() === PHP_SESSION_ACTIVE) {
                $user = $db->query("SELECT * FROM user_profile WHERE email=:email", array("email"=>$_SESSION['email']))[0];
                ?>
                <ul class="nav navbar-nav navbar-right">
                    <li class="dropdown">
                        <a href="#" class="dropdown-toggle nav-link" data-toggle="dropdown"><b><?=$_SESSION['email']?></b><span class="caret"></span></a>
                        <ul id="login-dp" class="dropdown-menu" style="min-width:325px;left:-200px;background-color:#343a40;border-color:#343a40;">
                            <li>
                                <div class="row">
                                    <div class="col-md-12">
					<ul class="list-group">
						<li class="list-group-item">Name: <span style="float:right"><?=$user['name'];?></span></li>
						<li class="list-group-item">Email: <span style="float:right"><?=$user['email'];?></span></li>
						<li class="list-group-item">Ethnicity: <span style="float:right"><?=$user['ethnicity'];?></span></li>
						<li class="list-group-item">Gender(s): <?php $gs=explode("\n", $user['gender']);foreach($gs as $g){ echo '<span style="float:right">'.$g.'</span><br>';} ?></li>
					</ul>
                                    </div>
                                    <div class="bottom text-center" style="border-color:#343a40;background-color:#343a40">
                                        <ul class="list-group">
						<li class="list-group-item"><a style="margin-right:5px;" href="logout.php"><b>Logout</b></a><?=$_SESSION['name']?></li>
					</ul>
                                    </div>
                                </div>
                            </li>
                        </ul>
                    </li>
                </ul>
                <?php
            } else {
                ?>
            <ul class="nav navbar-nav navbar-right">
                <li class="dropdown">
                    <a href="#" class="dropdown-toggle nav-link" data-toggle="dropdown"><b>Login</b> <span class="caret"></span></a>
                    <ul id="login-dp" class="dropdown-menu" style="min-width:325px;left:-250px;background-color:#343a40;border-color:#343a40;color:white">
                        <li>
                            <div class="row">
                                <div class="col-md-12">
                                    <form class="form" role="form" method="post" action="login.php" accept-charset="UTF-8" id="login-nav">
                                        <div class="form-group">
                                            <label class="sr-only" for="mail">Email address</label>
                                            <input type="email" class="form-control" name="mail" id="mail" placeholder="Email address" required>
                                        </div>
                                        <div class="form-group">
                                            <label class="sr-only" for="pwd">Password</label>
                                            <input type="password" class="form-control" id="pwd" name="pwd" placeholder="Password" required>
                                            <div class="help-block text-right" style="margin-top:5px;"><a href="" onclick="alert('So sad, you should remember it next time!');">Forgot the password?</a></div>
                                        </div>
                                        <div class="form-group">
                                            <button type="submit" class="btn btn-primary btn-block">Login</button>
                                        </div>
                                        <div class="checkbox">
                                            <label>
                                                <input name="remember" type="checkbox"> Remember me, I'm hungry!
                                            </label>
                                        </div>
                                    </form>
                                </div>
                                <div class="bottom text-center" style="border-color:#343a40;background-color:#343a40">
                                    New here? <a href="registration.php"><b>Join Us</b></a>
                                </div>
                            </div>
                        </li>
                    </ul>
                </li>
            </ul>
            <?php
            }
            ?>
        </div>
    </div>
</nav>
<div style="margin-top:3%"></div>
