#!/bin/bash

xx=$(netstat -avn|awk '{print $1"\t"$4"\t"$5"\t"$9}'|grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}")

# https://stackoverflow.com/questions/6212219/passing-parameters-to-a-bash-function
fnGetPidArgCmd () {
  while read -r foo;
  do
    xx1=$(ps -p `echo ${foo}|awk '{print $1}'|grep -Eo '([0-9]+)'` -o args|tail -n 1);
    echo  "${foo} ${xx1}";
  done
}

fnGetPidArgCmd4id () {
  while read -r foo;
  do
    xx1=$(ps -p ${foo} -o args|tail -n 1);
    echo  "${foo} ${xx1}";
  done
}


# fnGetPidArgCmd1 () {
#   while read -r foo;
#   do
#     xx1=$(ps -p `echo ${foo}|awk '{print $4}'` -o args|tail -n 1);
#     echo  "${foo}   ${xx1}";
#   done
# }

if [[ "f" = $1 ]];
  then
    echo -e "${xx}"|awk '{print $4" "$3}'|grep -Eo '([0-9]{2,} ([0-9]{1,3}\.){3}[0-9]{1,3})'|grep -Ev "127.0.0.1"|sort -u|uniq|fnGetPidArgCmd;
else
  #  echo -e "${xx}"|sort -u|uniq|grep -Ev "\*\.\*"|fnGetPidArgCmd;
  # pid, cmd args
  echo -e "${xx}"|awk '{print $4}'|grep -Eo '[0-9]{2,}'|sort -u|uniq|grep -Ev "\*\.\*"|fnGetPidArgCmd4id;
fi
