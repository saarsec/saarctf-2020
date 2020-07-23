<?php

class Menu
{    
    var $monday = array(1, 2, 3, 4);
    var $tuesday = array(5, 6, 7, 8);
    var $wednesday = array(9, 1, 2, 3);
    var $thursday = array(4, 9, 1, 7);
    var $friday = array(2, 5, 6, 3);

    public $day_count = 0;
	var $func = array('save'=>'store_menu');
	public function __construct($monday, $tuesday, $wednesday, $thursday, $friday){
		$this->monday=$monday;
		$this->tuesday=$tuesday;
		$this->wednesday=$wednesday;
		$this->thursday=$thursday;
		$this->friday=$friday;
	}
}
$max=27;
$monday=array(rand(1,$max),rand(1,$max),rand(1,$max),rand(1,$max));
$tuesday=array(rand(1,$max),rand(1,$max),rand(1,$max),rand(1,$max));
$wednesday=array(rand(1,$max),rand(1,$max),rand(1,$max),rand(1,$max));
$thursday=array(rand(1,$max),rand(1,$max),rand(1,$max),rand(1,$max));
$friday=array(rand(1,$max),rand(1,$max),rand(1,$max),rand(1,$max));
$m = new Menu($monday,$tuesday,$wednesday,$thursday,$friday);
echo serialize($m);

?>

