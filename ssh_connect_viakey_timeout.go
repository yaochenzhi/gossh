// golang naming: https://stackoverflow.com/questions/38616687/which-way-to-name-a-function-in-go-camelcase-or-semi-camelcase
// golang comment style? do as you like
// do you like?!
package main 

import "golang.org/x/crypto/ssh"
import "io/ioutil"
import "strings"
import "bytes"
import "time"
import "log"
import "net"
import "fmt"
import "os"


/*
  output:
    ''
    Failed to dial: dial tcp a.a.a.a:22: i/o timeout
    Failed to dial: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
*/
func main() {

    var (
        user = "root"
        port = "22"
        key = "/root/.ssh/id_rsa"
        timeout = 30 * time.Second
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
        // Timeout:  // not that valid sometimes. we do not use this param here.
    }

    // connection
    hostPort := fmt.Sprintf("%s:%s", host, port)
    connection, err := SSHDialTimeout("tcp", hostPort, sshConfig, timeout)
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

type Conn struct {
    net.Conn
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

func (c *Conn) Read(b []byte) (int, error) {
    err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
    if err != nil {
        return 0, err
    }
    return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
    err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
    if err != nil {
        return 0, err
    }
    return c.Conn.Write(b)
}

func SSHDialTimeout(network, addr string, config *ssh.ClientConfig, timeout time.Duration) (*ssh.Client, error) {
    conn, err := net.DialTimeout(network, addr, timeout)
    if err != nil {
        return nil, err
    }

    timeoutConn := &Conn{conn, timeout, timeout}
    c, chans, reqs, err := ssh.NewClientConn(timeoutConn, addr, config)
    if err != nil {
        return nil, err
    }
    client := ssh.NewClient(c, chans, reqs)

    // this sends keepalive packets every 2 seconds
    // there's no useful response from these, so we can just abort if there's an error
    go func() {
        t := time.NewTicker(2 * time.Second)
        defer t.Stop()
        for range t.C {
            _, _, err := client.Conn.SendRequest("keepalive@golang.org", true, nil)
            if err != nil {
                return
            }
        }
    }()
    return client, nil
}
