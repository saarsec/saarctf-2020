<?php
$page = 'reservation';
require_once "header.php";

function isReserved($row, $table, $seatrow, $seatnumber) {
    global $db;
    $data = array("row"=>$row, "tablenumber"=>$table, "seatrow"=>$seatrow, "seatnumber"=>$seatnumber);
    $ret = $db->query("SELECT * FROM seat WHERE row=:row AND tablenumber=:tablenumber AND seatrow=:seatrow AND seatnumber=:seatnumber;", $data);
    foreach ($ret as $y) {
        return $y['reserved_by'];
    }
    return false;
}

$own = null;
if (isset($_SESSION['email'])) {
    $mail = $_SESSION['email'];
    $own = $db->query("SELECT * FROM seat WHERE reserved_by=:mail", array("mail"=>$mail));
}

?>

<div class="container">
    <div class="row">
        <div>
            <h1 class="my-4">Seat Reservation:</h1>
            <?php if ($own) {echo '<h5>You have already reserved Seat '.$own[0]['row'].'.'.$own[0]['tablenumber'].'.'.$own[0]['seatrow'].'.'.$own[0]['seatnumber'].'</h5>';} ?>
        </div>
        <div id="notification_success" class="alert alert-success" style="display:none">
            <strong>Success!</strong> You have successfully reserved a seat!
        </div>
        <div id="notification_fail" class="alert alert-danger" style="display:none">
            <strong>Fail!</strong> This Seat is already taken or you already have reserved a seat!
        </div>
        <div style="margin-top:2%;margin-left:10%;width:100%;">
            <?php
            for ($a = 1; $a <= 3; $a++) {
                ?>
                <div class="row" style="margin-bottom:5%">
                    <?php
                    for ($b = 1; $b <= 4; $b++) {
                        ?>
                        <div class="col-3">
                            <div class="row">
                                <div class="col-3" style="width:30px;margin-top:3px;margin-right:-25px;">
                                    <?php
                                    for ($i = 1; $i <= 8; $i++) {
                                        ?>
                                        <div class="row" style="cursor:pointer;margin-top: 5px;">
                                            <img width="25" height="25" onclick="$('.popover').popover('hide');$(this.nextElementSibling).popover('show');" src="pictures/chair.png" class="rotate270"/><a href="#" title="Seat <?=$a?>.<?=$b?>.1.<?=$i?>" data-toggle="popover" <?php $x=isReserved($a, $b,1, $i);if ($x){echo "data-content='Status: Reserved by<br> ".htmlspecialchars($x)."'";} ?> data-trigger="focus"></a>
                                        </div>
                                        <?php
                                    }
                                    ?>
                                </div>
                                <div class="col-3" style="width:42px;margin-right:20px;">
                                    <img width="250" height="42" src="pictures/table.png" class="rotate90" style="margin-left:-105px;margin-top:105px;"/>
                                </div>
                                <div class="col-3" style="width:30px;margin-top:3px;">
                                    <?php
                                    for ($i = 1; $i <= 8; $i++) {
                                        ?>
                                        <div class="row" style="cursor:pointer;margin-top: 5px;">
                                            <img width="25" height="25" onclick="$('.popover').popover('hide');$(this.nextElementSibling).popover('show');" src="pictures/chair.png" class="rotate90"/><a href="#" title="Seat <?=$a?>.<?=$b?>.2.<?=$i?>" data-toggle="popover" <?php $x=isReserved($a, $b,2, $i);if ($x){echo "data-content='Status: Reserved by<br> ".htmlspecialchars($x)."'";} ?>data-trigger="focus"></a>
                                        </div>
                                        <?php
                                    }
                                    ?>
                                </div>
                            </div>
                        </div>
                        <?php
                    }
                    ?>
                </div>
                <?php
            }
            ?>
        </div>
        <div class="row">
            <h1 class="my-4">Picnic Places: <button type="button" class="btn" style="margin-left: 20px;margin-bottom:10px;margin-top:5px;" onclick="reserveSeat({'innerText': 'Picnic 0.0.0.0'})">Get Yours!</button></h1>
        </div>
        <div class="row" style="width:100%;background-color:#53c953c7;">
            <?php
            $picnics = $db->query("SELECT * FROM seat WHERE row=0 AND tablenumber=0 AND seatrow=0 AND reserved_by!='';");
            foreach ($picnics as $p) {
                echo '<div class="col-1" style="margin: 5px;">';
                echo '<img width="75" height="75" onclick="$(\'.popover\').popover(\'hide\');$(this.nextElementSibling).popover(\'show\');" src="pictures/picnic.png"><a href="#" title="Seat 0.0.0.'.$p['seatnumber'].'" data-toggle="popover" data-content="Status: Reserved by<br>'.htmlspecialchars($p['reserved_by']).'" data-trigger="focus"></a>';
                echo '</div>';
            }
            ?>
        </div>

    </div>
</div>
<div id="popupdiv" style="display: none">
    Status: free to reserve<br>
    <button type="button" class="btn" style="float:right;margin-bottom:10px;margin-top:5px;" onclick="reserveSeat(this.parentElement.previousElementSibling)">Take it!</button>
</div>
<script>

    $(document).ready(function() {
        $('[data-toggle="popover"]').popover({
            html : true,
            content: function() {
                console.log(this);
                return $('#popupdiv').html();
            }
        });
    });

    function reserveSeat(elm) {
        var seat_id = elm.innerText.split(' ')[1];

        var mail = "<?= $_SESSION['email']; ?>";

        var tmp= seat_id.split('.');

        var data = {
            'user': mail,
            'row': tmp[0],
            'table': tmp[1],
            'seatrow': tmp[2],
            'seatnumber': tmp[3]
        };

        reserve(data);
    }

    function reserve(data) {

        var xhr = new XMLHttpRequest();

        var url = "do_reservation.php";

        var params = [];
        for(var k in data) {
            if (data.hasOwnProperty(k)) {
                params.push(encodeURIComponent(k) + "=" + encodeURIComponent(data[k]));
            }
        }
        params = params.join("&");

        xhr.open("POST", url, true);

        xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

        xhr.onreadystatechange = function() {
            if(xhr.readyState === 4) {
                if (xhr.status === 200) {
                    var x = document.getElementById("notification_success");
                    x.style.display = 'block';
                    setTimeout(function() {
                        document.location.reload();
                    }, 1000);
                } else {
                    var y = document.getElementById("notification_fail");
                    y.style.display = 'block';
                    setTimeout(function() {
                        y.style.display = 'none';
                    }, 2500);
                }
            }
        };

        xhr.send(params);
    }
</script>
<?php
require_once "footer.html";
?>
