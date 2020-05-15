package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	/*var dataPool = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}*/

	scanner := bufio.NewScanner(file)
	foundUsers := ""
	seenBrowsers := []string{}
	var user User
	i := 0
	//user := dataPool.Get().(*User)

	for scanner.Scan() {
		lines := scanner.Bytes()

		uniqueBrowsers := 0
		i++

		if !bytes.Contains(lines, []byte("Android")) && !bytes.Contains(lines, []byte("MSIE")) {
			continue
		}
		err := json.Unmarshal(lines, &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {

			if ok := strings.Contains(browserRaw, "Android"); ok {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}
		for _, browserRaw := range user.Browsers {

			if ok := strings.Contains(browserRaw, "MSIE"); ok {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					//log.Printf("SLOW New browser: %s, first seen: %s", browserRaw, user.Name)
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}
		if !(isAndroid && isMSIE) {
			continue
		}
		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i-1, user.Name, email)
	}
	//dataPool.Put(user)
	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}