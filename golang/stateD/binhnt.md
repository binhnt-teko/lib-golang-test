# Test process in D state  
- Compile 
cc -o dstate dstate.c
- Run 
./dstate & sleep 0.1; ps -o pid,state,cmd -p "$!"
 pkill -P "$!"


cc -o dstate1 dstate1.c

./dstate1 & sleep 0.1; ps -o pid,state,cmd -p "$!"
pkill -P "$!"
ps -ef | grep dstate1


cc -o signal-bitmap signal-bitmap.c
