# File Integrity

Compare o conteúdo de uma pasta com uma coleta anterior. Cada arquivo regular recebe um SHA-256 calculado em fluxo. A comparação informa arquivos adicionados, alterados, removidos e caminhos inconclusivos. Uma pasta ilegível torna inconclusivos os arquivos conhecidos abaixo dela.

## Windows · PowerShell

Instale Go 1.22 ou superior. Na pasta do repositório, execute:

```powershell
go build -o file-integrity.exe .
mkdir amostra
Set-Content 'amostra/ação local.txt' 'primeira versão'
./file-integrity.exe coletar amostra base.json
Set-Content 'amostra/ação local.txt' 'segunda versão'
./file-integrity.exe verificar base.json amostra
```

## Linux · Bash

```bash
go build -o file-integrity .
mkdir amostra
printf 'primeira versão' > 'amostra/ação local.txt'
./file-integrity coletar amostra base.json
printf 'segunda versão' > 'amostra/ação local.txt'
./file-integrity verificar base.json amostra
./file-integrity coletar amostra atual.json
./file-integrity comparar base.json atual.json
```

Guarde os manifestos fora da pasta analisada, em local confiável. A coleta exige um novo arquivo de saída. A comparação escreve uma lista JSON ordenada no terminal. Código 0 indica sucesso sem divergência, 1 indica mudanças e 2 indica erro ou coleta parcial. Consulte o campo falhas do manifesto para o motivo por caminho.

A ferramenta ignora links simbólicos, dispositivos e diretórios vazios. Não fotografa o sistema de arquivos atomicamente. A verificação de estabilidade usa identidade, tamanho e horário de modificação e não protege contra adulteração deliberada dos metadados ou do manifesto.

## Validação

Execute `go test ./...` e `go vet ./...`. O workflow inclui Windows e Linux. Sem dependências além da biblioteca padrão.

Licença MIT.
