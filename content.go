package restapi

import (
	"io/ioutil"
	"net/http"
	"github.com/yiyedata/restapi/utils"

	"github.com/pquerna/ffjson/ffjson"
)

type H map[string]interface{}
type KV map[string]string

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	D          H
	handlers   []HandlerFunc // middleware
	index      int
	kvs        map[string]string
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		D:      make(H, 0),
		index:  -1,
		kvs:    make(KV, 0),
	}
}
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	if c.index < s {
		c.handlers[c.index](c)
	}
}
func (c *Context) GetPara(key string) string {
	v, _ := c.kvs[key]
	return v
}
func (c *Context) SetPara(key string, value string) {
	c.kvs[key] = value
}
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Body() *utils.JsonMap {
	bodyBytes, _ := ioutil.ReadAll(c.Req.Body)
	return utils.Json2map(bodyBytes)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, message []byte) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write(message)
}

func (c *Context) JSON(obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(200)
	buf, err := ffjson.Marshal(&obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	} else {
		c.Writer.Write(buf)
	}
	// encoder := json.NewEncoder(c.Writer)
	// if err := encoder.Encode(obj); err != nil {
	// 	http.Error(c.Writer, err.Error(), 500)
	// }
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
