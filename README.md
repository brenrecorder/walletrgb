README

CLI Wallet for RGB in golang currently testing server is up and set already
5 coins for free (BETA)

Compile (windows,linux,termux):<Br>
go mod init wallet<br>
go mod tidy<br>
go build wallet.go
edit server.cnf or set -server 127.0.0.1:8080

server included (windows exe file)
usage: coinserver -password X -difficulty X 

miner included
usage:
miner.exe -server 127.0.0.1:8080 to start mining blockchain

<img src="https://github.com/brenrecorder/walletrgb/blob/main/screenwallet.png?raw=true"></img>
<img src="https://github.com/brenrecorder/walletrgb/blob/main/screenserver.png?raw=true"></img>
<img src="https://github.com/brenrecorder/walletrgb/blob/main/screenminer.png?raw=true"></img>
