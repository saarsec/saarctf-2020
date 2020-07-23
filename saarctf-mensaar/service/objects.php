<?php

class Menu
{
    var $monday = array();
    var $tuesday = array();
    var $wednesday = array();
    var $thursday = array();
    var $friday = array();

    public $day_count = 0;
    var $func = array('save'=>'store_menu');
}

class Feedback
{
    var $rating = 0;
    var $comment = '';

    var $func = array('evaluate'=>'ignore_feedback');

    public function ignore_feedback($feedback)
    {
        print_r('Rating: '.$feedback->rating);
        print_r('Comment: '.$feedback->comment);
    }

    public function __wakeup()
    {
        @$this->func['evaluate']($this);
    }
}
