scp combigram.linux james@carbocation.com:/data/bin/upload.combigram.linux

ssh -n -f james@carbocation.com "sh -c 'cd /data/bin/; killall combigram.linux; mv upload.combigram.linux combigram.linux;GOMAXPROCS=8  nohup ./combigram.linux > /dev/null 2>&1 &'"
