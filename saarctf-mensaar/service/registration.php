<?php

if( isset($_POST['name']) && isset($_POST['email']) && isset($_POST['pwd']) ) {

    require_once "dblib.php";

    $db = new Db;
    $db->init();

    $token = hash_hmac ( 'sha3-256' , $_POST['name'], $_POST['pwd']);

    if ( isset($_POST['gender']) )
        $gender = implode(",", $_POST['gender']);
    else
        $gender = "";

    if ( isset($_POST["ethnicity"]) )
        $ethnicity = $_POST["ethnicity"];
    else
        $ethnicity = "";

	$data = array("name" => $_POST['name'], "email" => $_POST['email'], "pwd" => password_hash($_POST['pwd'],PASSWORD_DEFAULT), "gender" => $gender, "ethnicity" => $ethnicity);
	try{
		$db->execute("INSERT INTO user_profile VALUES (:name, :email, :pwd, null, :gender, :ethnicity)", $data);
	}
	catch(Exception $e) {
		echo $e;
		echo 'This user already exists : (';
		http_response_code(200);
		die();
	}
    session_start();
    $_SESSION["name"] = $_POST['name'];
    $_SESSION["email"] = $_POST['email'];
    header("Location: index.php");
    die();

} else {

    $page = 'home';
    require_once "header.php";
    ?>
    <div class="container">
        <div class="row">
            <h1 class="my-4">Registration:</h1>
            <div style="margin-top:2%;width:100%;">
                <form action="registration.php" method="POST">
                    <div class="form-group">
                        <label for="name">Username:</label>
                        <input type="text" class="form-control" name="name" id="name">
                    </div>
                    <div class="form-group">
                        <label for="email">Email address:</label>
                        <input type="email" class="form-control" name="email" id="email">
                    </div>
                    <div class="form-group">
                        <label for="reg_pwd">Password:</label>
                        <input type="password" class="form-control" name="pwd" id="reg_pwd">
                    </div>
                    <div class="form-group">
                        <label for="gender">Your Gender(s):</label>
                        <select multiple class="form-control" name ="gender[]" id="gender">
                            <option>Abimegender</option>
                            <option>Absorgender</option>
                            <option>Adamasgender</option>
                            <option>Adeptogender</option>
                            <option>Aerogender</option>
                            <option>Aesthetgender</option>
                            <option>Aethergender</option>
                            <option>Affectugender</option>
                            <option>Agender</option>
                            <option>Agenderfluid</option>
                            <option>Agenderflux</option>
                            <option>Alexigender</option>
                            <option>Aliusgender</option>
                            <option>Ambigender</option>
                            <option>Amaregender</option>
                            <option>Ambonec</option>
                            <option>Amicagender</option>
                            <option>Amogender</option>
                            <option>Amorgender</option>
                            <option>Androgyne</option>
                            <option>Anesigender</option>
                            <option>Angeligender</option>
                            <option>Angenital</option>
                            <option>Anogender</option>
                            <option>Anongender</option>
                            <option>Antegender</option>
                            <option>Antigender</option>
                            <option>Anxiegender</option>
                            <option>Anvisgender</option>
                            <option>Apagender</option>
                            <option>Apconsugender</option>
                            <option>Apogender</option>
                            <option>Apollogender</option>
                            <option>Aporagender</option>
                            <option>Aptugender</option>
                            <option>Aquarigender</option>
                            <option>Archaigender</option>
                            <option>Arifluid</option>
                            <option>Arigender</option>
                            <option>Arithmagender</option>
                            <option>Argogender</option>
                            <option>Astergender</option>
                            <option>Astralgender</option>
                            <option>Atmosgender</option>
                            <option>Autigender</option>
                            <option>Autogender</option>
                            <option>Axigender</option>
                            <option>Batgender</option>
                            <option>Bigender</option>
                            <option>Bigenderfluid</option>
                            <option>Biogender</option>
                            <option>Blizzgender</option>
                            <option>Boggender</option>
                            <option>Bordergender</option>
                            <option>Borderfluid</option>
                            <option>Boyflux</option>
                            <option>Brevigender</option>
                            <option>Burstgender</option>
                            <option>Cadensgender</option>
                            <option>Cadogender</option>
                            <option>Caedogender</option>
                            <option>Caelgender</option>
                            <option>Cancegender</option>
                            <option>Canisgender</option>
                            <option>Caprigender</option>
                            <option>Carmigender</option>
                            <option>Cassflux</option>
                            <option>Cassgender</option>
                            <option>Caveagender</option>
                            <option>Cavusgender</option>
                            <option>Cendgender</option>
                            <option>Cennedgender</option>
                            <option>Ceterofluid</option>
                            <option>Ceterogender</option>
                            <option>Chaosgender</option>
                            <option>Cheiragender</option>
                            <option>Circgender</option>
                            <option>Cloudgender</option>
                            <option>Cocoongender</option>
                            <option>Cogitofluid</option>
                            <option>Cogitogender</option>
                            <option>Coigender</option>
                            <option>Collgender</option>
                            <option>Colorgender</option>
                            <option>Comgender</option>
                            <option>Commogender</option>
                            <option>Condigender</option>
                            <option>Contigender</option>
                            <option>Corugender</option>
                            <option>Cosmicgender</option>
                            <option>Cryptogender</option>
                            <option>Crystagender</option>
                            <option>Cyclogender</option>
                            <option>Daimogender</option>
                            <option>Deaboy</option>
                            <option>Delphigender</option>
                            <option>Demifluid/flux</option>
                            <option>Demigender</option>
                            <option>Digigender</option>
                            <option>Diurnalgender</option>
                            <option>Domgender</option>
                            <option>Drakefluid</option>
                            <option>Dryagender</option>
                            <option>Dulcigender</option>
                            <option>Duragender</option>
                            <option>Eafluid</option>
                            <option>Earthgender</option>
                            <option>Egender</option>
                            <option>Ectogender</option>
                            <option>Effreu</option>
                            <option>Egogender</option>
                            <option>Ekragender</option>
                            <option>Eldrigender</option>
                            <option>Elegender</option>
                            <option>Elementgender</option>
                            <option>Elissogender</option>
                            <option>Enbyfluid</option>
                            <option>Endogender</option>
                            <option>Energender</option>
                            <option>Entheogender</option>
                            <option>Entrogender</option>
                            <option>Equigender</option>
                            <option>Espigender</option>
                            <option>Evaisgender</option>
                            <option>Exgender</option>
                            <option>Exiccogender</option>
                            <option>Existigender</option>
                            <option>Expecgender</option>
                            <option>Explorogender</option>
                            <option>Faegender</option>
                            <option>Fascigender</option>
                            <option>Faunagender</option>
                            <option>Fawngender</option>
                            <option>Felisgender</option>
                            <option>Femfluid</option>
                            <option>Femgender</option>
                            <option>Firegender</option>
                            <option>Fissgender</option>
                            <option>Flirtgender</option>
                            <option>Flowergender</option>
                            <option>Fluidflux</option>
                            <option>Foggender</option>
                            <option>Frostgender</option>
                            <option>Fuzzgender</option>
                            <option>Gemelgender</option>
                            <option>Gemigender</option>
                            <option>Geminigender</option>
                            <option>Genderale</option>
                            <option>Genderamas</option>
                            <option>Genderblank</option>
                            <option>Genderblur</option>
                            <option>Gendercosm</option>
                            <option>Genderdormant</option>
                            <option>Gendereaux</option>
                            <option>Genderflora</option>
                            <option>Genderflight</option>
                            <option>Genderflow</option>
                            <option>Genderfluid</option>
                            <option>Genderflux</option>
                            <option>Genderfuzz</option>
                            <option>Gendermaverick</option>
                            <option>Gendernegative</option>
                            <option>Gender-Neutral</option>
                            <option>Genderplasma</option>
                            <option>Genderpositive</option>
                            <option>Genderpunk</option>
                            <option>Genderqueer</option>
                            <option>Gendersea</option>
                            <option>Genderstrange</option>
                            <option>Gendervague</option>
                            <option>Gendervex</option>
                            <option>Gendervoid</option>
                            <option>Genderwitched</option>
                            <option>Gendfleur</option>
                            <option>Girlflux</option>
                            <option>Glassgender</option>
                            <option>Glimragender</option>
                            <option>Glitchgender</option>
                            <option>Gossagender</option>
                            <option>Greengender</option>
                            <option>Greygender</option>
                            <option>Gyraboy</option>
                            <option>Gyragender</option>
                            <option>Gyragirl</option>
                            <option>Healegender</option>
                            <option>Heliogender</option>
                            <option>Hemigender</option>
                            <option>Horogender</option>
                            <option>Hydrogender</option>
                            <option>Hypogender</option>
                            <option>Illusogender</option>
                            <option>Impediogender</option>
                            <option>Imperigender</option>
                            <option>Inersgender</option>
                            <option>Intergender</option>
                            <option>Invisigender</option>
                            <option>Iragender</option>
                            <option>Jupitergender</option>
                            <option>Juxera</option>
                            <option>Kingender</option>
                            <option>Kynigender</option>
                            <option>Lamingender</option>
                            <option>Leogender</option>
                            <option>Lethargender</option>
                            <option>Leukogender</option>
                            <option>Levigender</option>
                            <option>Liberique</option>
                            <option>Libragender</option>
                            <option>Librafluid</option>
                            <option>Lichtgender</option>
                            <option>Lipsigender</option>
                            <option>Locugender</option>
                            <option>Lovegender</option>
                            <option>Ludogender</option>
                            <option>Lysigender</option>
                            <option>Magigender</option>
                            <option>Maringender</option>
                            <option>Marfluid</option>
                            <option>Mascfluid</option>
                            <option>Mascugender</option>
                            <option>Maverique</option>
                            <option>Medeigender</option>
                            <option>Melogender</option>
                            <option>Mirrorgender</option>
                            <option>Molligender</option>
                            <option>Moongender</option>
                            <option>Mosaigender</option>
                            <option>Musicgender</option>
                            <option>Mutaregender</option>
                            <option>Mutogender</option>
                            <option>Mystigender</option>
                            <option>Nanogender</option>
                            <option>Narkissigender</option>
                            <option>Necrogender</option>
                            <option>Nesciogender</option>
                            <option>Neurogender</option>
                            <option>Neutrois</option>
                            <option>Nobifluid</option>
                            <option>Nocturnalgender</option>
                            <option>Non-binary</option>
                            <option>Novigender</option>
                            <option>Nubilagender</option>
                            <option>Nullgender</option>
                            <option>Nyctogender</option>
                            <option>Obruogender</option>
                            <option>Offgender</option>
                            <option>Omnigay</option>
                            <option>Orbgender</option>
                            <option>Owlgender</option>
                            <option>Paragender</option>
                            <option>Pendogender</option>
                            <option>Perigender</option>
                            <option>Perogender</option>
                            <option>Personagender</option>
                            <option>Perospike</option>
                            <option>Pictogender</option>
                            <option>Pixelgender</option>
                            <option>Polygender</option>
                            <option>Polygenderflux</option>
                            <option>Portiogender</option>
                            <option>Praegender</option>
                            <option>Preciogender</option>
                            <option>Preterbinary</option>
                            <option>Primusgender</option>
                            <option>Privagender</option>
                            <option>Proxvir</option>
                            <option>Quivergender</option>
                            <option>Quoigender</option>
                            <option>Salugender</option>
                            <option>Schrodigender</option>
                            <option>Scorigender</option>
                            <option>Scorpifluid</option>
                            <option>Scorpigender</option>
                            <option>Seagender</option>
                            <option>Selenogender</option>
                            <option>Sequigender</option>
                            <option>Shellgender</option>
                            <option>Skygender</option>
                            <option>Spesgender</option>
                            <option>Spikegender</option>
                            <option>Stargender</option>
                            <option>Staticgender</option>
                            <option>Stratogender</option>
                            <option>Subgender</option>
                            <option>Subfluid</option>
                            <option>Surgender</option>
                            <option>Swampgender</option>
                            <option>Sychnogender</option>
                            <option>Systemfluid</option>
                            <option>Systemgender</option>
                            <option>Tachigender</option>
                            <option>Tangender</option>
                            <option>Tauragender</option>
                            <option>Technogender</option>
                            <option>Telegender</option>
                            <option>Tempgender</option>
                            <option>Temporagender</option>
                            <option>Tenuigender</option>
                            <option>Tragender</option>
                            <option>Traumatgender</option>
                            <option>Trigender</option>
                            <option>Turbogender</option>
                            <option>Ungender</option>
                            <option>Vaguefluid</option>
                            <option>Vagueflux</option>
                            <option>Vaguegender</option>
                            <option>Vapogender</option>
                            <option>Vectorgender</option>
                            <option>Veloxigender</option>
                            <option>Venngender</option>
                            <option>Venufluid</option>
                            <option>Verangender</option>
                            <option>Vestigender</option>
                            <option>Vibragender</option>
                            <option>Videgender</option>
                            <option>Videogender</option>
                            <option>Virgender</option>
                            <option>Vocigender</option>
                            <option>Voidfluid</option>
                            <option>Voidgender</option>
                            <option>Witchgender</option>
                            <option>Xenogender</option>
                            <option>Xirl</option>
                            <option>Xoy</option>
                            <option>Xumgender</option>
                            <option>Zodiacgender</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="ethnicity">Your ethnicity:</label>
                        <select class="form-control" name = "ethnicity" id="ethnicity">
                            <option>White privileged Male</option>
                            <option>Suppressed Female</option>
                            <option>Trash redneck</option>
                            <option>Eco-hippie</option>
                            <option>Other</option>
                        </select>
                    </div>
                    <button type="submit" class="btn btn-primary">Submit</button>
                </form>
            </div>
        </div>
    </div>

    <?php
    require_once "footer.html";
}
?>
