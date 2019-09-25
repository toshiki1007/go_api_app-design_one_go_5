package main

import (
		"github.com/gin-gonic/gin"
		"math/rand"
		"time"
		"net/http"
		"bufio"
		"bytes"
		"io"
    "strings"
    "github.com/mattn/go-encoding"
    "golang.org/x/net/html"
    "golang.org/x/net/html/charset"
)
func text(resp *http.Response) (string, error) {
	br := bufio.NewReader(resp.Body)
	var r io.Reader = br
	if data, err := br.Peek(1024); err == nil {
			if _, name, ok := charset.DetermineEncoding(data, resp.Header.Get("content-type")); ok {
					if enc := encoding.GetEncoding(name); enc != nil {
							r = enc.NewDecoder().Reader(br)
					}
			}
	}

	var buffer bytes.Buffer
	doc, err := html.Parse(r)
	if err != nil {
			return "", err
	}
	walk(doc, &buffer)
	return buffer.String(), nil
}
func walk(node *html.Node, buff *bytes.Buffer) {
	if node.Type == html.TextNode {
			data := strings.Trim(node.Data, "\r\n ")
			if data != "" {
					buff.WriteString("\n")
					buff.WriteString(data)
			}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
			switch strings.ToLower(node.Data) {
			case "script", "style", "title":
					continue
			}
			walk(c, buff)
	}
}
func main() {
	r := gin.Default()
	rand.Seed(time.Now().UnixNano())

	r.GET("/num", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"num": rand.Intn(6)+1,
		})
	})

  url := "https://qiita.com/api/v2/items"

	r.GET("/qiita", func(c *gin.Context) {
		resp, _ := http.Get(url)

		defer resp.Body.Close()
		s, _ := text(resp)
		c.JSON(200, gin.H{
			"qiita": s,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
