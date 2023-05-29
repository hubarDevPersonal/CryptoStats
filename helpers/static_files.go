package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/shurcooL/httpgzip"
	"github.com/twmb/murmur3"
)

type StaticFiles struct {
	apiFS http.FileSystem
}

func NewStaticFiles(apiFS http.FileSystem) *StaticFiles {
	return &StaticFiles{
		apiFS: apiFS,
	}
}

func (sf *StaticFiles) AttachTo(engine *echo.Echo, prefix string) {
	_ = vfsutil.WalkFiles(sf.apiFS, "/", sf.addHandler(sf.apiFS, prefix, engine))

	// Add redirect for non-slash terminated API urls (because Swagger in a browser won't work without the slash on the end)
	noSlash := strings.TrimSuffix(prefix, "/")
	withSlash := noSlash + "/"
	engine.GET(noSlash, func(c echo.Context) error {
		return c.Redirect(302, withSlash)
	})
}

func (sf *StaticFiles) addHandler(fs http.FileSystem, prefix string, engine *echo.Echo) func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
	return func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if fi.IsDir() {
			return nil
		}
		handler := sf.NewFileHandler(fs, path)

		path = prefix + path

		// handle the file as specified
		engine.GET(path, handler)

		if !strings.HasSuffix(path, "index.html") {
			return nil
		}

		// handle index.html when the directory is requested with trailing slash
		dirPath := path[:len(path)-len("index.html")]
		engine.GET(dirPath, handler)

		if len(dirPath) == 1 {
			return nil
		}

		// handle index.html when the directory is requested without trailing slash
		dirPath = dirPath[:len(dirPath)-1]
		engine.GET(dirPath, handler)

		return nil
	}
}

func (sf *StaticFiles) NewFileHandler(fs http.FileSystem, path string) echo.HandlerFunc {
	contentType := ""
	cacheControl := "no-cache, max-age=0"
	if strings.HasSuffix(path, ".js") {
		cacheControl = "public, max-age=3600"
		contentType = "application/javascript"
	} else if strings.Contains(path, "favicon.ico") {
		cacheControl = "public, max-age=3600"
	}

	etag := sf.GenerateEtag(fs, path)

	return func(c echo.Context) error {
		if etagIn := c.Request().Header.Get("If-None-Match"); etagIn != "" && etag != "" {
			if etag == etagIn {
				return c.NoContent(http.StatusNotModified)
			}
		}

		f, _ := fs.Open(path)
		defer f.Close()

		hdr := c.Response().Header()
		hdr.Set("Cache-Control", cacheControl)
		hdr.Set("ETag", etag)
		if contentType != "" {
			hdr.Set("ContentParams-Type", contentType)
		}

		httpgzip.ServeContent(c.Response(), c.Request(), path, time.Time{}, f)
		return nil
	}
}

func (sf *StaticFiles) GenerateEtag(fs http.FileSystem, path string) string {
	b, err := vfsutil.ReadFile(fs, path)
	if err != nil {
		return ""
	}

	h1, h2 := murmur3.Sum128(b)
	return fmt.Sprintf("%X%X", h1, h2)
}
