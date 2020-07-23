<?php
$page='home';
require_once "header.php";
?>

<div class="container">
    <div id="main" class="row">
        <div class="col-lg-3">
            <h1 class="my-4">Menus:</h1>
            <div class="list-group" id="side-navbar">
                <a href="#" class="list-group-item active" style="color:white;font-weight:bold;" onclick="activate(this);">Monday</a>
                <a href="#" class="list-group-item" style="color:black;" onclick="activate(this);">Tuesday</a>
                <a href="#" class="list-group-item" style="color:black;" onclick="activate(this);">Wednesday</a>
                <a href="#" class="list-group-item" style="color:black;" onclick="activate(this);">Thursday</a>
                <a href="#" class="list-group-item" style="color:black;" onclick="activate(this);">Friday</a>
            </div>
        </div>
        <!-- /.col-lg-3 -->
        <?php
            $days = array("monday", "tuesday", "wednesday", "thursday", "friday");
            foreach ($days as $day) {
                if ($day === "monday") {
                    ?>
                        <div id="menu_<?=$day?>" class="col-lg-9" style="display: block">
                    <?php
                } else {
                    ?>
                            <div id="menu_<?=$day?>" class="col-lg-9" style="display: none">
                    <?php
                }

            ?>
                <h1 class="my-4">&nbsp;</h1>
                <ul class="list-group">
                    <?php
                    $foods = $db->query("SELECT * FROM menu INNER JOIN food ON menu.food=food.id AND menu.day=:day ORDER BY menu.date LIMIT 5;", array("day"=>$day));

                    foreach ($foods as $food) {
                        ?>
                        <li class="list-group-item">
                            <div class="row">
                                <div class="col-2">
                                    <img src="<?= $food['pic_src'] ?>"
                                         style="margin-left:-7px;max-width:128px;max-height:128px;">
                                </div>
                                <div class="col-10">
                                    <h5><?= $food['name'] ?></h5>
                                    <p style="margin-bottom:0px"><?= $food['ingredients'] ?></p>
                                    <span class="text-warning" style="margin-right:5px;">
                                        <?php
                                            for ($i = 1; $i <= 5; $i++) {
                                                if ($i <= $food['rating']) {
                                                    echo "&#9733; ";
                                                } else {
                                                    echo "&#9734; ";
                                                }
                                            }
                                        ?>
                                    </span> <?= $food['rating'] ?> stars
                                </div>
                            </div>
                        </li>
                        <br>
                        <?php
					}
				?>
                </ul>
            </div>
            <?php
            }
        ?>
    </div>
</div>
<script>

    var activate = function(elem) {
        var items = document.getElementById('side-navbar').children;
        var menus = document.getElementById('main').children;
        var day = elem.innerText.toLowerCase();
        for (var i = 0; !!items[i]; i++) {
            menus[i+1].style.display = (!(menus[i+1].id.indexOf(day)>0) ? 'none' : 'block');
            items[i].classList.remove("active");
            items[i].style.color = 'black';
            items[i].style.fontWeight = 'normal';
        }
        elem.classList.add("active");
        elem.style.color = 'white';
        elem.style.fontWeight = 'bold';
    }

</script>
<?php
require_once "footer.html";
?>
