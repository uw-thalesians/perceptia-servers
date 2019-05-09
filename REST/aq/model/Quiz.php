<?php

class Quiz
{
    private $_id;
    private $_keyword;
    private $_summary;

    private $_image;

    private $_questions;
    private $_timestamp;
    private $_source;

    public function __construct($init)
    {
        $this->_id = $init['id'];
        $this->_keyword = $init['keyword'];
        $this->_summary = $init['summary'];
        $this->_image = $init['image'];
        $this->_timestamp = $init['when'];
        $this->_source = $init['source'];

        $this->_questions = array();
    }

    public function __get($var)
    {
        $tmp = '_' . $var;

        if(property_exists('Quiz', $tmp))
        {
            return $this->$tmp;
        }

        return null;

    }

    public function __set($var, $val)
    {
        $tmp = '_' . $var;

        $rw_props = array("_questions");

        if(in_array($tmp, $rw_props))
        {
            $this->$tmp = $val;
        }
    }

}

?>
