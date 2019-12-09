try:
    import commands
except Exception:
    import subprocess as commands


def ssh_check(host):

    conn_timeout = 30

    jump_server = ""
    jump_server_info = {
        "port": 22,
        "username": "username",
        "password": "password",
        "timeout": conn_timeout
    }
    jump_cmd = "sshc {}".format(host)

    with paramiko.SSHClient() as ssh:
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(jump_server, **jump_server_info)
        stdin, stdout, stderr = ssh.exec_command(jump_cmd, timeout=conn_timeout)
        output = stdout.read()

    if output == "":
        code, status = 0, "ok"
    elif "none publickey" in output:
        code, status = 1, "none publickey"
    else:
        code, status = 2, "unok"
    return code, status


if __name__ == "__main__":
    output = commands.getoutput("sshc host")
    if output == "":
        code, status = 0, "ok"
    elif "none publickey" in output:
        code, status = 1, "none publickey"
    else:
        code, status = 2, "unok"