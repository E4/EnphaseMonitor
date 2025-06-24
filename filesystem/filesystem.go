package filesystem

import (
  "net/http"
  "strings"
)

type RestrictedFileSystem struct {
  FileSystem http.FileSystem
  AllowSource bool
}

func (rfs RestrictedFileSystem) Open(path string) (http.File, error) {
  if strings.HasPrefix(path,"/src") && rfs.AllowSource==false { return rfs.BadRequest() }
  f, err := rfs.FileSystem.Open(path)
  if err != nil { return rfs.BadRequest() }

  s, err := f.Stat()
  if err != nil { return rfs.BadRequest() }

  if !s.IsDir() { return f, nil }

  index := strings.TrimSuffix(path, "/") + "/index.html"
  if _, err := rfs.FileSystem.Open(index); err != nil { return rfs.BadRequest() }

  return f, nil
}

func (rfs RestrictedFileSystem) BadRequest() (http.File, error) {
  f, err := rfs.FileSystem.Open("404.html")
  if err != nil { return nil, err }
  return f, nil
}