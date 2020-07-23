<?php
$page='feedback';
require_once "header.php";
require_once "objects.php";

?>
<div class="container">
    <div class="row">
        <?php
        if( isset($_POST['food']) && isset($_POST['rating']) && isset($_SESSION['email']) )
        {
            $user_mail = $_SESSION['email'];
            $token = hash( 'sha3-256' , $_SESSION['email']);
            $f = new Feedback();
            $f->rating = $_POST['rating'];
            $f->comment = $_POST['comment'];
            $data = array("user"=>$user_mail, "food"=>$_POST['food'], "obj"=>serialize($f), "cook_token"=>$token);
            $db->execute("INSERT INTO feedback VALUES (:food, :user, :obj, :cook_token)", $data);
            ?>
            <div style="width: 50%;margin: 0 auto;">
                <h1 class="my-4">Thanks for your Feedback!</h1>
                <p>Please note that your mensa does not give a single fuck about it.</p>
            </div>
            <?php
        } else {
        ?>
        <h1 class="my-4">Send us Feedback:</h1>
        <div style="margin-top:2%;width:100%;">
            <form action="feedback.php" method="POST">
                <div class="form-group">
                    <label for="food">Select Food:</label>
                    <select class="form-control" id="food" name="food">
                        <?php

                        $foods = $db->query("SELECT id, name FROM food;");

                        foreach ($foods as $food) {
                            ?>
                            <option value="<?= $food['id'] ?>"><?= $food['name'] ?></option>
                            <?php
                        }
                        ?>
                    </select>
                </div>
                <div class="form-group">
                    <label for="rating">Rating:</label><br>
                    <span class="text-warning" style="font-size:2em;">
                        <?php
                        for ($i = 1; $i <= 5; $i++) {
                            ?>
                            <span id="star<?= $i ?>" style="cursor:pointer;"
                                  onclick="setStar(<?= $i ?>);">&#9733;</span>
                            <?php
                        }
                        ?>
                    </span>
                    <input type="number" max="5" min="1" class="form-control" id="rating" name="rating" value=5 hidden>
                </div>
                <div class="form-group">
                    <label for="comment">Comment:</label>
                    <textarea class="form-control" rows="5" id="comment" name="comment"></textarea>
                </div>
                <button type="submit" class="btn btn-primary">Submit</button>
            </form>
        </div>
    </div>
</div>
    <script>
        function setStar(x) {
            for (var i = 1; i <= 5; i++)
                $('#star' + i).html('&#9734;');
            for (var j = 1; j <= x; j++)
                $('#star' + j).html('&#9733;');
            $('#rateing').val(x);
        }
    </script>
<?php
}
require_once "footer.html";
?>
