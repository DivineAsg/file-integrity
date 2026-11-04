package main

import (
 "os"
 "path/filepath"
 "testing"
)

func TestCiclo(t *testing.T){
 root:=t.TempDir()
 file:=filepath.Join(root,"ação local")
 write:=func(s string){t.Helper();if err:=os.WriteFile(file,[]byte(s),0600);err!=nil{t.Fatal(err)}}
 write("")
 base:=filepath.Join(t.TempDir(),"base.json")
 check:=func(want int,args ...string){t.Helper();got,err:=run(args);if got!=want{t.Fatalf("%v: código %d, erro %v",args,got,err)}}
 check(0,"coletar",root,base)
 check(2,"coletar",root,base)
 check(2,"coletar",root,filepath.Join(root,"manifesto"))
 check(0,"verificar",base,root)
 write("mudou")
 check(1,"verificar",base,root)
 if err:=os.Remove(file);err!=nil{t.Fatal(err)}
 check(1,"verificar",base,root)
 check(0,"comparar",base,base)
}

func TestPastaInacessivel(t *testing.T){
 a:=Manifest{1,map[string]string{"p/a":"x","q":"y"},map[string]string{}}
 b:=Manifest{1,map[string]string{},map[string]string{"p":"acesso negado"}}
 got:=compare(a,b)
 if len(got)!=3||got[1]!=(Change{"p/a","inconclusivo"})||got[2]!=(Change{"q","removido"}){t.Fatal(got)}
}
