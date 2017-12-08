pid=`ps -ef | grep albums_web | grep -v grep | awk '{print $2}'`
kill -9 $pid