package main

import (
 "path"
 "sort"
)

type Change struct { Path string `json:"caminho"`; State string `json:"estado"` }

func unknown(m Manifest,key string)bool{
 for {
  if _,ok:=m.Errors[key];ok{return true}
  if key=="."{return false}
  key=path.Dir(key)
 }
}

func compare(a,b Manifest)[]Change{
 keys:=map[string]bool{}
 for _,m:=range []Manifest{a,b}{
  for k:=range m.Files{keys[k]=true}
  for k:=range m.Errors{keys[k]=true}
 }
 changes:=[]Change{}
 for k:=range keys{
  old,inOld:=a.Files[k]
  now,inNew:=b.Files[k]
  state:=""
  switch{
  case unknown(a,k)||unknown(b,k):state="inconclusivo"
  case !inOld:state="adicionado"
  case !inNew:state="removido"
  case old!=now:state="alterado"
  }
  if state!=""{changes=append(changes,Change{k,state})}
 }
 sort.Slice(changes,func(i,j int)bool{return changes[i].Path<changes[j].Path})
 return changes
}
