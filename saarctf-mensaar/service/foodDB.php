<?php
$page='foodDB';
require_once "header.php";
?>
<script src="https://cdn.datatables.net/1.10.16/js/jquery.dataTables.min.js"></script>
<script src="https://cdn.datatables.net/1.10.16/js/dataTables.bootstrap4.min.js"></script>

<div class="container">
    <div class="row">
        <h1 class="my-4">Food Database Browser:</h1>
        <div style="margin-top:1%;width:100%;">
            <form class="form-inline">
            <table id="food_table" class="table table-striped table-bordered" cellspacing="0" width="100%" style="margin-top: 2%"><!--class="table table-hover">-->
                <thead>
                <tr>
                    <th>Food</th>
                    <th>Ingredients</th>
                    <th>Rating</th>
                    <th>With Maggi?</th>
                </tr>
                </thead>
                <tbody>
                    <?php
                    $foods = $db->query("SELECT * FROM food;", null);

                    foreach ($foods as $food) {
                        ?>
                        <tr>
                            <td><?=$food['name']?></td>
                            <td><?=$food['ingredients']?></td>
                            <td><?php
                                for ($i = 1; $i <= 5; $i++) {
                                    if ($i <= $food['rating']) {
                                        echo "&#9733; ";
                                    } else {
                                        echo "&#9734; ";
                                    }
                                }
                                ?></td>
                            <td>Yes</td>
                        </tr>
                        <?php
                    }
                    ?>
                </tbody>
            </table>
            </form>
        </div>
    </div>
</div>
<script>

    function align_table_stuff() {
        var bot = document.getElementsByClassName("pagination");
        bot[0].style.float = "right";
        bot[0].onclick = function(){align_table_stuff();};
        var top = document.getElementsByClassName("col-sm-12 col-md-6");
        if (top[1].innerHTML.indexOf("Show") === -1)
            top[1].parentNode.insertBefore(top[1], top[0]);
        top[0].children[0].children[0].style.float = "left";
        top[1].children[0].children[0].style.float = "right";
        top[0].children[0].children[0].style.width = "66%";
        top[0].children[0].children[0].style.fontSize = '0em';
        top[0].children[0].children[0].children[0].placeholder = "Search...";
        top[0].children[0].children[0].children[0].style.width = "100%";
        top[1].children[0].children[0].children[0].onchange = function(){bot[0].style.float = "right";};
        top[1].children[0].children[0].children[0].style.margin = "0px 5px";
        top[0].children[0].children[0].children[0].onkeyup = function(){align_table_stuff();};
    }

    $(document).ready(function() {
        $('#food_table').DataTable();
        align_table_stuff();
    });

</script>

<?php
require_once "footer.html";
?>
