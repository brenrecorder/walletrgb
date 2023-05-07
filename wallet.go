package main

import (
    "fmt"
    //"log"
    "net/http"
    "encoding/hex"
    "crypto/md5"
    "os"
    "io"
    "strings"
    "strconv"
	"time"
	"gitlab.com/david_mbuvi/go_asterisks"
	"bytes"
	//"flag"
	"github.com/nexidian/gocliselect"
    )
    var walletaddress string
    var passwordmd5 string
    var usedpassword string 
    
    var flagpassword string
    var blockchainchecked bool = false
   var server string = "86.84.201.181:95"
    func main() {
		
menu := gocliselect.NewMenu("Main menu")
menu.AddItem("Make transaction", "1")
menu.AddItem("Refresh balance", "2")
menu.AddItem("History", "3")
menu.AddItem("Receive", "4")
menu.AddItem("Other server", "5")
menu.AddItem("Offline coins", "6")
menu.AddItem("Exit", "7")

		var menuchoice int 
setserver(false)
		if isFileExisting("mywallet.db") {
		if ReadUserSettings() {
			fmt.Print("\033[H\033[2J")
			if (RetrieveAmountWallet() <0) {
				fmt.Print("Wallet not found..\n\n") 
				setserver(true)
				fmt.Print("\033[H\033[2J")
			}
			strAmount := fmt.Sprintf("%.3f",RetrieveAmountWallet())

		//fmt.Print("RGB Wallet\nAddress:\t"+walletaddress+"\nAmount:\t\t"+strAmount+"\n\n1:Make Transaction\n2:Refresh balance\n3:Other server\n4:Offline coins\n5:Exit\n\n")
		//fmt.Scanln(&menuchoice)
		var blockchainvalid string = "invalid"
		if (blockchainchecked)  {blockchainvalid = "valid"} 
		fmt.Print("RGB Wallet\nBlockchain:\t"+blockchainvalid+"\n\nAddress:\t"+walletaddress+"\nAmount:\t\t"+strAmount+"\n\n")
		menuchoice, _ = strconv.Atoi(menu.Display())
		if menuchoice >0 {  } else { main()}
		if menuchoice == 1 { 
			MakeTransaction() 
			time.Sleep(3 * time.Second)
			main()
			}
		if menuchoice == 2 { 
			main()
			}
			if menuchoice == 3 {
				TransActionHistory(walletaddress)
				fmt.Println("\npress any key to exit..")
				fmt.Scanln()
				main()
			}
			if menuchoice == 4 {
			ReceiveCoins(walletaddress)
			main()
			}
			if menuchoice == 5 { 
				setserver(true)
				main()
			}
			if menuchoice == 6 {
				fmt.Print("\033[H\033[2J")
				fmt.Print("RGB Wallet\n\nAddress:\t"+walletaddress+"\nAmount:\t\t"+strAmount+"\n\n")
						
				menuofflinecoins := gocliselect.NewMenu("Offline coins")
				menuofflinecoins.AddItem("Import coins", "1")
				menuofflinecoins.AddItem("Export coins", "2")
				menuofflinecoins.AddItem("Exit", "3")
				var menuchoiceb int
				menuchoiceb, _ = strconv.Atoi(menuofflinecoins.Display())
				
				if menuchoiceb == 1 {
				var coincode string
				fmt.Print("Enter coin code: ")
				fmt.Scanln(&coincode)
					fmt.Println(ImportCoins(coincode))
				fmt.Println("press any key to continue..")
				fmt.Scanln()
					main()
					}
				if menuchoiceb == 2 {
				var amount float64
				fmt.Print("Amount coins to export: ")
				fmt.Scanln(&amount)
				
				coinamountexp := fmt.Sprintf("%v", amount)
				fmt.Print("\nCoins: \t" + coinamountexp + "\n")
				fmt.Println("Code: \t" +GetCoinsOffline(amount))
				 fmt.Println("Store the code somewhere then press any key to continue..")
				fmt.Scanln()
				main()
					}
				if menuchoiceb == 3 { main() }


			}
		if menuchoice == 7 { 
			os.Exit(0)
			}
	} else {
		fmt.Print("RGB Wallet\nPassword invalid..\n")
	}
	} else {
		fmt.Print("RGB Wallet\nEnter new wallet password: ")
		//var newwalletpassword string
		//fmt.Scanln(&newwalletpassword)
		newwalletpasswordb, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	
		if err != nil {
			// handle error
		}
		newwalletpassword := bytes.NewBuffer(newwalletpasswordb).String()
		resp, err := http.Get("http://" + server + "/coinserver?action=createwallet&password=" + newwalletpassword)
		if err != nil {
			// handle error
		}

		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		if err != nil {
		fmt.Print("RGB Wallet\nServer offline..")
		} else {
		fmt.Println(string(body))
		 walletaddress := strings.Split(string(body), ":")
		 CreateWalletFile(walletaddress[1],newwalletpassword)
		 main()
	}
	}
	}
func setserver(setnewserver bool) bool {
			if isFileExisting("server.cnf") && setnewserver == false{
				b, err := os.ReadFile("server.cnf") // just pass the file name
				if err != nil {
				fmt.Print(err)
				} else {
					server = string(b)
					return true
				}
		}
	fmt.Print("\nRGB Wallet\nServer: "+server+"\nEnter new server address: ")
	var setserver string
	fmt.Scanln(&setserver)
	

    if len(setserver) >0 {
	f, err := os.Create("server.cnf")
    if err != nil {
        fmt.Println(err)
        return false
    }
    _, err = f.WriteString(setserver)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return false
    }

    server = setserver
}
    return true
}
func ReceiveCoins(wallet string) {
var receiveAmount float64
fmt.Print("Enter payment amount: ")
fmt.Scanln(&receiveAmount)
var boolinhistory bool = false
var boolreceived bool = false
fmt.Println("Awaiting payment..")

		resp, err := http.Get("http://" + server + "/coinserver?action=transactionhistory")
		if err != nil {
			fmt.Println("Server offline or no internet connection..")	
		}
		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		oldtransactions := strings.Split(string(body), "\n")
		var oldTransActionsBlockChain = map[int]string{}
		for i, block := range oldtransactions {
		histblocksplit := strings.Split(string(block), ":")
			oldTransActionsBlockChain[i] = histblocksplit[0]
		}
		
for boolreceived == false {
		resp, err := http.Get("http://" + server + "/coinserver?action=transactionhistory")
		if err != nil {
			
			fmt.Println("Server offline or no internet connection..")
			
		}
//bchash:amount:from:to
		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		if err != nil { fmt.Println("Error") } else {
		alltransaction := strings.Split(string(body), "\n")
		
		//fmt.Println("\nTransaction history\nFrom\t\t\tTo\t\tAmount\n")
		for _, block := range alltransaction {
		splitblock := strings.Split(block, ":")
		boolinhistory = false
		for _, blockold := range oldTransActionsBlockChain {
			if len(splitblock) > 2 && blockold == splitblock[0] { boolinhistory = true }
		}
			//fmt.Println(block)
			
			if len(splitblock) > 2 && boolinhistory == false  {
			strAmounttoReceive := fmt.Sprintf("%.3f",receiveAmount)
			if (splitblock[2] == wallet || splitblock[3] == wallet) && splitblock[1] == strAmounttoReceive  {
			fmt.Println("Payment done, " + strAmounttoReceive + " From: " + splitblock[2])
			boolreceived = true
			//fmt.Println(splitblock[2] + "\t" + splitblock[3] + "\t\t" + splitblock[1] + "")
			}
			}
		}

		}
time.Sleep(3 * time.Second)
}

}
func TransActionHistory(wallet string) {
		resp, err := http.Get("http://" + server + "/coinserver?action=transactionhistory")
		if err != nil {
			
			fmt.Println("Server offline or no internet connection..")
			
		}
//bchash:amount:from:to
		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		if err != nil { fmt.Println("Error") } else {
		alltransaction := strings.Split(string(body), "\n")
		
		fmt.Println("\nTransaction history\nFrom\t\t\tTo\t\tAmount\n")
		for _, block := range alltransaction {
			//fmt.Println(block)
			splitblock := strings.Split(block, ":")
			if len(splitblock) > 2 {
			if splitblock[2] == wallet || splitblock[3] == wallet {
			fmt.Println(splitblock[2] + "\t" + splitblock[3] + "\t\t" + splitblock[1] + "")
			}
			}
		}

		}
}

func GetCoinsOffline(amount float64) string {

	coinamount := fmt.Sprintf("%v", amount)
		resp, err := http.Get("http://" + server + "/coinserver?action=makeofflinecoin&wallet=" + walletaddress + "&password=" +usedpassword+"&amount="+coinamount)
		if err != nil {
			return "failed"
			fmt.Println("Server offline or no internet connection..")
			
		}

		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		if err != nil { } else {
			retrvcoincoide := strings.Split(string(body), ":")
			return string(retrvcoincoide[1])
		}
		return "failed"
}
func ImportCoins(coincode string) string {

		resp, err := http.Get("http://" + server + "/coinserver?action=importofflinecoins&wallet=" + walletaddress + "&coincode=" +coincode +"&password="+usedpassword)
		if err != nil {
			return "failed"
			fmt.Println("Server offline or no internet connection..")
			
		}
		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		if err != nil { } else {
			//imported := strings.Split(string(body), ":")
			return string(body)
		}
		return "failed"
}

func RetrieveAmountWallet() float64 {
		resp, err := http.Get("http://" + server + "/coinserver?action=getamount&wallet=" + walletaddress + "&password=" +usedpassword)
		if err != nil {
			return -1
			fmt.Println("Server offline or no internet connection..")
			
		}

		defer resp.Body.Close()
		body, err:= io.ReadAll(resp.Body)
		if err != nil {
		fmt.Print("RGB Wallet\nServer offline..")
		} else {
			walletamount := strings.Split(string(body), ":")
			walletamountfloat, _ := strconv.ParseFloat(walletamount[1], 3)
			return walletamountfloat
		}
		return 0.0
}
func MakeTransaction() {
	var adrTo string
	var amount float64
	
	fmt.Print("\nNew Transaction\nAddress: ")
	fmt.Scanln(&adrTo)
	fmt.Print("Amount: ")
	fmt.Scanln(&amount)
	s := fmt.Sprintf("%v", amount)
	resp, err := http.Get("http://" + server + "/coinserver?action=transaction&password="+usedpassword+"&from="+walletaddress+"&to="+adrTo+"&amount="+s)
		if err != nil {
			
		}
	defer resp.Body.Close()
	body, err:= io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("\n" + string(body))
		} else {
			fmt.Println("\n" + string(body))
		}
}
func ReadUserSettings() bool {
	if len(usedpassword) > 0 {
		return true
	}
	b, err := os.ReadFile("mywallet.db") // just pass the file name
    if err != nil {
        fmt.Print(err)
    } else {
		var checkpassword string
		 readvarsfile := strings.Split(string(b), ":")
		 walletaddress = readvarsfile[0]
		 passwordmd5 = readvarsfile[1]
		 fmt.Print("RGB Wallet\n\n")
		 if len(flagpassword) <1 {
	
		fmt.Print("\nAddress:\t"+walletaddress+"\n\nLogin with wallet password\nPassword: ")

			newwalletpasswordb, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	
		if err != nil {
			// handle error
		}
		checkpassword = bytes.NewBuffer(newwalletpasswordb).String()
	} else {
		checkpassword = flagpassword
	}
		
		if stringtoMD5(checkpassword) == passwordmd5 {
			usedpassword = checkpassword
					 if CheckBlockChainValid() {blockchainchecked = true} else {blockchainchecked = false}
			return true
		}
	}
	return false
}
func CreateWalletFile(address string, password string) bool {
	f, err := os.Create("mywallet.db")
    if err != nil {
        fmt.Println(err)
        return false
    }
    _, err = f.WriteString(address + ":" +stringtoMD5(password))
    if err != nil {
        fmt.Println(err)
        f.Close()
        return false
    } else {
		fmt.Println("Wallet: " + address + " password: " + password + " created")
		return true
	}
}

func CheckBlockChainValid() bool {
		resp, err := http.Get("http://" + server + "/coinserver?action=blockchain")
		if err != nil {
			
		}
	defer resp.Body.Close()
	body, err:= io.ReadAll(resp.Body)
		if err != nil {
			return false
		} else {
			response := strings.Split(string(body), ":")
			if response[1] == "true" { 
				return true
			} else {
				return false
			}
		}
		return false
}
func isFileExisting(filename string) bool {
   info, err := os.Stat(filename)
   if os.IsNotExist(err) {
      return false
   }
   return !info.IsDir()
}

func stringtoMD5(text string) string {
   hash := md5.Sum([]byte(text))
   return hex.EncodeToString(hash[:])
}
