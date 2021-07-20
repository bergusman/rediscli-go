package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/bergusman/rediscli-go/resp"
)

func run(host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(conn)
	}

	prompt := conn.RemoteAddr().String()

	cmds := bufio.NewReader(os.Stdin)
	dec := resp.NewDecoder(conn)

	for {
		fmt.Print(prompt + "> ")
		cmd, err := cmds.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if cmd == "exit\n" {
			os.Exit(0)
		}

		_, err = conn.Write([]byte(cmd))
		if err != nil {
			log.Fatal(err)
		}

		res, err := dec.Decode()
		if err != nil {
			log.Fatal(err)
		}

		print(res)

		switch v := res.(type) {
		case []interface{}:
			if len(v) > 0 {
				if vv, ok := v[0].([]byte); ok && string(vv) == "subscribe" {
					fmt.Println("Reading messages... (press Ctrl-C to quit)")
					for {
						res, err := dec.Decode()
						if err != nil {
							log.Fatal(err)
						}
						print(res)
					}
				}
			}
		}
	}
}

func RESPString(v interface{}) string {
	w := &strings.Builder{}

	var write func(v interface{})
	write = func(v interface{}) {
		switch v := v.(type) {
		default:
			fmt.Fprint(w, v)
		case nil:
			w.WriteString("nil")
		case string:
			fmt.Fprintf(w, "%q", v)
		case []byte:
			fmt.Fprintf(w, "%q", v)
		case []interface{}:
			w.WriteString("[")
			for i, vv := range v {
				write(vv)
				if i != len(v)-1 {
					w.WriteString(" ")
				}
			}
			w.WriteString("]")
		}
	}

	write(v)
	return w.String()
}

func print(v interface{}) {
	fmt.Println(RESPString(v))
}

func main() {
	host := ":6379"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}
	run(host)
}
