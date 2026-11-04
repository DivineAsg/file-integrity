package main

import (
 "encoding/hex"
 "encoding/json"
 "fmt"
 "io/fs"
 "os"
 "path/filepath"
 "strings"
)

func load(name string)(Manifest,error){
 var m Manifest
 data,err:=os.ReadFile(name)
 if err!=nil{return m,err}
 if err=json.Unmarshal(data,&m);err!=nil{return m,err}
 if m.Version!=1||m.Files==nil||m.Errors==nil{return m,fmt.Errorf("manifesto inválido")}
 for key,sum:=range m.Files{
  _,err:=hex.DecodeString(sum)
  if !fs.ValidPath(key)||key=="."||len(sum)!=64||err!=nil{return m,fmt.Errorf("entrada inválida: %q",key)}
 }
 for key:=range m.Errors{if !fs.ValidPath(key){return m,fmt.Errorf("caminho de falha inválido")}}
 return m,nil
}

func run(args []string)(int,error){
 if len(args)!=3{return 2,fmt.Errorf("uso: file-integrity coletar PASTA NOVO.json | verificar BASE.json PASTA | comparar BASE.json ATUAL.json")}
 if args[0]=="coletar"{
  root,err:=filepath.EvalSymlinks(args[1]);if err!=nil{return 2,err}
  parent,err:=filepath.EvalSymlinks(filepath.Dir(args[2]));if err!=nil{return 2,err}
  root,err=filepath.Abs(root);if err!=nil{return 2,err}
  dest,err:=filepath.Abs(filepath.Join(parent,filepath.Base(args[2])));if err!=nil{return 2,err}
  rel,err:=filepath.Rel(root,dest)
  if err==nil&&rel!=".."&&!strings.HasPrefix(rel,".."+string(filepath.Separator)){return 2,fmt.Errorf("salve o manifesto fora da pasta analisada")}
  m,err:=scan(args[1]);if err!=nil{return 2,err}
  f,err:=os.OpenFile(dest,os.O_CREATE|os.O_EXCL|os.O_WRONLY,0600);if err!=nil{return 2,err}
  err=json.NewEncoder(f).Encode(m);closeErr:=f.Close()
  if err!=nil||closeErr!=nil{os.Remove(dest);return 2,fmt.Errorf("falha ao gravar manifesto: %v %v",err,closeErr)}
  if len(m.Errors)>0{return 2,fmt.Errorf("coleta parcial: consulte falhas no manifesto")}
  return 0,nil
 }
 a,err:=load(args[1]);if err!=nil{return 2,err}
 var b Manifest
 switch args[0]{case "verificar":b,err=scan(args[2]);case "comparar":b,err=load(args[2]);default:return 2,fmt.Errorf("operação desconhecida")}
 if err!=nil{return 2,err}
 changes:=compare(a,b)
 if err=json.NewEncoder(os.Stdout).Encode(changes);err!=nil{return 2,err}
 if len(a.Errors)+len(b.Errors)>0{return 2,nil}
 if len(changes)>0{return 1,nil}
 return 0,nil
}

func main(){code,err:=run(os.Args[1:]);if err!=nil{fmt.Fprintln(os.Stderr,err)};os.Exit(code)}
