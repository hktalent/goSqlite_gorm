#!/bin/bash
domain=$1

# amass 执行，后期考虑集成一个golang版本
# 分布式执行命令

#  将命令发送到后台执行
tmux new -s subdomain -d
tmux send -t "subdomain" "ssh -i ~/.ssh/id_rsa -p $newSshPort -C root@${newIp} \"/root/tools/doSubdomain4Server.sh ${domain}\nexit\n\"" Enter

