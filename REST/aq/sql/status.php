<?php
abstract class STATUS {

    const __default = self::RECVD;
    
    const REQ_RECVD = 0;
    const RETRVD_MEDIA = 1;
    const MEDIA_SPLIT = 2;
    const QG_BEGINS = 3;
    const QG_COMPLETED = 4;
    const READY = 4;
    const __COUNT = self::READY+1;

    public static function getCount() { return self::__COUNT; }

    const STATUS_STRINGS = array("We got your request!",
                             "We've retrieved your media!",
                             "Content Analysis completed!",
                             "Performing Question Generation!",
                             "Your Quiz is ready!");

    const STATUS_NOTFOUND = "Requested resource not found!";
};
?>