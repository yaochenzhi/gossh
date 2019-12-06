// golang naming: https://stackoverflow.com/questions/38616687/which-way-to-name-a-function-in-go-camelcase-or-semi-camelcase
// golang comment style? do as you like
// do you like?!
package main 

import "golang.org/x/crypto/ssh"
import "io/ioutil"
import "strings"
import "bytes"
import "log"
import "net"
import "fmt"
import "os"


func main() {

    var (
        user = "root"
        port = "22"
        key = "/root/.ssh/id_rsa"
    )

    host, command := parseArgs()

    sshConfig := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            PublicKeyFile(key),
        },
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
    }

    // connection
    hostPort := fmt.Sprintf("%s:%s", host, port)
    connection, err := ssh.Dial("tcp", hostPort, sshConfig)
    if err != nil {
        log.Fatal(fmt.Errorf("Failed to dial: %s", err))
    }
    defer connection.Close()

    // session
    if command != "" {
        session, err := connection.NewSession()
        if err != nil {
            log.Fatal(fmt.Errorf("Failed to create session: %s", err))
        }
        defer session.Close()

        var stdoutBuf bytes.Buffer
        session.Stdout = &stdoutBuf

        session.Run(command)
        fmt.Println(session.Stdout)
    }

}

func PublicKeyFile(file string) ssh.AuthMethod {
    buffer, err := ioutil.ReadFile(file)
    if err != nil {
        return nil
    }

    key, err := ssh.ParsePrivateKey(buffer)
    if err != nil {
        return nil
    }
    return ssh.PublicKeys(key)
}

func parseArgs()(string, string) {
    host, command := "", ""
    if len(os.Args) > 1 {
        host = os.Args[1]
        if len(os.Args) > 2{
            command = strings.Join(os.Args[2:], " ")
        }
    }else{
        echoUsage()
    }
    return host, command
}

func echoUsage() {
    _exe := os.Args[0]
    log.Fatal(fmt.Sprintf("%s host [command]", _exe))
}
