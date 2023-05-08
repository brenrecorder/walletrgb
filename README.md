README

CLI Wallet for RGB in golang currently testing server is up and set already
5 coins for free (BETA)

Compile (windows,linux,termux):<Br>
go mod init wallet<br>
go mod tidy<br>
go build wallet.go
edit server.cnf or set -server 127.0.0.1:8080

server included (windows exe file)<br>
usage: coinserver -password X -difficulty X 

miner included<br>
Please mine 5 minutes or more before using<br>
usage:<br>
miner.exe -server 127.0.0.1:8080 to start mining blockchain

<img src="https://github.com/brenrecorder/walletrgb/blob/main/pictures/screenwallet.png?raw=true"></img>
<img src="https://github.com/brenrecorder/walletrgb/blob/main/pictures/screenserver.png?raw=true"></img>
<img src="https://github.com/brenrecorder/walletrgb/blob/main/pictures/screenminer.png?raw=true"></img>
