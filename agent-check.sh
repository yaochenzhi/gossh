#!/usr/bin/env bash
# * * * * * agent-check.sh
# After that:
#     "source /data/host/hang/agent.log; sshc(ssh_connect) host [cmd]"

_dir="/data/host/hang"

ssh_private_keys="
    /root/.ssh/id_rsa
"

agent_sock="${_dir}/agent.sock"
agent_log="${_dir}/agent.log"
agent_cmd="/usr/bin/ssh-agent -a $agent_sock"


function start_agent(){
    rm -f $agent_sock

    ${agent_cmd} >$agent_log
    sed -i '/echo.*/d' $agent_log

    eval `cat $agent_log`
    for key in $ssh_private_keys;
    do
        /usr/bin/ssh-add $key
    done
}

function main(){
    running=`ps -ef | grep "${agent_cmd}" | grep -v grep`
    if [ "$running" = "" ];then
        echo "Agent is NOT running"
        echo "Starting agent ..."
        start_agent
        echo "Agent started !"
    else
        echo "Agent is running"
        echo ${running}
    fi
}

main
# <<< END