<?php
$param=$_SERVER['argv'];

$param[0]="date:".date('Y-m-d H:i:s',time());

$log=implode(" ",$param);

echo $log;

file_put_contents("./fsnotify.log",$log.PHP_EOL,FILE_APPEND);