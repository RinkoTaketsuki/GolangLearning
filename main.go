package main // 表示当前文件含有 func main，可执行

import (
	"fmt"
	"net/url"
)

func main() {
	u, err := url.Parse("https://username:password@www.host.name:12345/path/t%20o/source/./../source?key1=value1&key2=value%202#title%201")
	if err != nil {
		panic(err)
	}
	scheme := u.Scheme
	username := u.User.Username()
	password, _ := u.User.Password()
	authentication := u.User.String()
	hostname := u.Hostname()
	port := u.Port()
	socketAddress := u.Host
	path := u.Path
	rawPath := u.EscapedPath()
	query := u.Query()
	rawQuery := u.RawQuery
	fragment := u.Fragment
	rawFragment := u.EscapedFragment()
	objects := []any{scheme, username, password, authentication, hostname, port, socketAddress, path, rawPath, query, rawQuery, fragment, rawFragment}
	for _, obj := range objects {
		fmt.Printf("%v\n", obj)
	}
	fmt.Printf("%v\n", u)
	fmt.Printf("%+v\n", u)
	fmt.Printf("%#v\n", u)
}
