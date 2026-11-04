package main

import (
 "crypto/sha256"
 "fmt"
 "io"
 "io/fs"
 "os"
 "path/filepath"
)

type Manifest struct {
 Version int `json:"versao"`
 Files map[string]string `json:"arquivos"`
 Errors map[string]string `json:"falhas"`
}

func digest(path string, initial fs.FileInfo) (string,error) {
 f,err:=os.Open(path)
 if err!=nil{return "",err}
 defer f.Close()
 before,err:=f.Stat()
 if err!=nil{return "",err}
 if !os.SameFile(initial,before)||!before.Mode().IsRegular(){return "",fmt.Errorf("arquivo substituído")}
 h:=sha256.New()
 n,err:=io.Copy(h,f)
 if err!=nil{return "",err}
 after,err:=os.Lstat(path)
 if err!=nil{return "",err}
 if !os.SameFile(before,after)||n!=before.Size()||after.Size()!=n||!before.ModTime().Equal(after.ModTime()){return "",fmt.Errorf("arquivo instável durante a leitura")}
 return fmt.Sprintf("%x",h.Sum(nil)),nil
}

func scan(root string)(Manifest,error){
 m:=Manifest{1,map[string]string{},map[string]string{}}
 info,err:=os.Lstat(root)
 if err!=nil{return m,err}
 if !info.IsDir(){return m,fmt.Errorf("a raiz precisa ser um diretório real")}
 err=filepath.WalkDir(root,func(path string,d fs.DirEntry,problem error)error{
  key,err:=filepath.Rel(root,path)
  if err!=nil{return err}
  key=filepath.ToSlash(key)
  if problem!=nil{m.Errors[key]=problem.Error();return nil}
  if !d.Type().IsRegular(){return nil}
  st,err:=d.Info()
  if err==nil {
   var sum string
   sum,err=digest(path,st)
   if err==nil{m.Files[key]=sum}
  }
  if err!=nil{m.Errors[key]=err.Error()}
  return nil
 })
 return m,err
}
