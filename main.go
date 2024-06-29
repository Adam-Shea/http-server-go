package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type HTTPCode string
type ContentType string

const (
	success        HTTPCode = "200 OK"
	error          HTTPCode = "500 Internal Server Error"
	notImplemented HTTPCode = "501 Not Implemented"
	lost           HTTPCode = "404 Not Found"
)

const (
	html ContentType = "Content-Type: text/html;"
	ico  ContentType = "Content-Type: image/vnd.microsoft.icon;"
	png  ContentType = "Content-Type: image/png;"
	jpg  ContentType = "Content-Type: image/jpeg;"
	svg  ContentType = "Content-Type: image/svg+xml;"
	webp ContentType = "Content-Type: image/webp;"
	gif  ContentType = "Content-Type: image/gif;"
	css  ContentType = "Content-Type: text/css;"
	js   ContentType = "Content-Type: text/javascript;"
	none ContentType = ""
)

func createHttpHeaders(errorCode HTTPCode, contentType ContentType) []byte {
	return []byte("HTTP/1.1 " + string(errorCode) + "\n" + string(contentType) + "\n\n")
}

func handleGet(conn net.Conn, path string) {
	b, err := os.ReadFile("serve/" + path)
	if err != nil {
		fmt.Printf("Error reading  : %#v\n", err)
		conn.Write(createHttpHeaders(lost, none))
		conn.Close()
	} else {
		var header = none
		if len(path) >= 5 && path[len(path)-5:] == ".html" {
			header = html
		} else if len(path) >= 4 && path[len(path)-4:] == ".htm" {
			header = html
		} else if len(path) >= 4 && path[len(path)-4:] == ".png" {
			header = png
		} else if len(path) >= 4 && path[len(path)-4:] == ".jpg" {
			header = jpg
		} else if len(path) >= 4 && path[len(path)-4:] == ".svg" {
			header = svg
		} else if len(path) >= 4 && path[len(path)-4:] == ".webp" {
			header = webp
		} else if len(path) >= 4 && path[len(path)-4:] == ".gif" {
			header = gif
		} else if len(path) >= 4 && path[len(path)-4:] == ".css" {
			header = css
		} else if len(path) >= 3 && path[len(path)-3:] == ".js" {
			header = js
		} else if len(path) >= 4 && path[len(path)-4:] == ".ico" {
			header = ico
		}
		conn.Write(createHttpHeaders(success, header))
		conn.Write(b)
	}
}

func main() {
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading   : %#v\n", err)
				conn.Write(createHttpHeaders(error, none))
				conn.Close()
			}

			var method string
			var path string

			fmt.Printf("Remote address : %v\n", conn.RemoteAddr())
			for i, requestLine := range strings.Split(string(buf[:len]), "\n") {
				if requestLine == "\r" {
					break
				} else if i == 0 {
					method = strings.Split(requestLine, " ")[0]
					path = strings.Split(requestLine, " ")[1]
					fmt.Printf("Path           : %v\n", path)
					if path == "" || path == "/" {
						path = "index.html"
					} else {
						if path[0] == '/' {
							path = path[1:]
						}
						for {
							if path[0] == '.' {
								path = path[1:]
							} else {
								break
							}
						}
					}
				} else {
					// parse headers
				}
			}

			switch method {
			case "GET":
				handleGet(conn, path)
				break
			case "POST":
				conn.Write(createHttpHeaders(notImplemented, none))
				break
			case "PUT":
				conn.Write(createHttpHeaders(notImplemented, none))
				break
			case "DELETE":
				conn.Write(createHttpHeaders(notImplemented, none))
				break
			default:
				conn.Write(createHttpHeaders(notImplemented, none))
				break
			}

			conn.Close()
		}(conn)
	}
}
