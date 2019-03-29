# ATHAN HOME
A fully customizable Athan for the Google Home Devices

### Installation

Athan Home Requires Golang to build the binary, and Apache or other WebServer + PHP + SQL
to fully function

Install these Golang dependencies.

```sh
$ go get -u github.com/evalphobia/google-home-client-go
$ go get -u github.com/gojektech/heimdall
$ go get -u github.com/micro/mdns
```
### Setting up the Google Homes
1. Open the Google Home App
2. At the top center click "Add"
3. Then click "Create Speaker Group"
4. Select all google homes you wish to play athan
5. Name the Speaker Group "Group Cast"
6. Make sure no other speaker group exists, if it does delete it, if you dont, the Athan may not work

### Setting up the binary
Setup for the most part is pretty simple given you have some basic knowledge of using linux or have atleasted used a raspberry.
For running this I recommend using a raspberry pi Zero W, its small and very usfull it will allow you to enjoy the full functionality of this neat program

Instaling Golang
We first must install GoLang (this isnt required when running the actuall binary but will be required when trying to build the binrary)

Please Follow this tutorial: 
```
https://tecadmin.net/install-go-on-ubuntu/
```

#### Building the binnary
First go to the Users folder in the panel and move Both "main.go" and "config.json" to there.
```sh
$ go build
```
Install Screen to use in background
```sh
$ sudo apt-install screen
```
Then run the screen command
```sh
$ screen
$ ./Users or main.go
$ Ctrl + A + D
```

### Setting up the panel
Once you have setup the Application move the Binary and Config.json file to the Users Folder
once you have done that run the command


Setting Up the WebServer
Please Install and Setup This first (You may use Apache):
```
https://mediatemple.net/community/products/developer/204405534/install-nginx-on-ubuntu
```
Then Install and Setup Phpmyadmin:
```
https://pimylifeup.com/raspberry-pi-mysql-phpmyadmin/
```
Once you have installed SQL and PHP My admin go to domain or ip/phpmyadmin ex. 127.0.0.1/phpmyadmin
create a new database, and click import and import athan.sql and make sure to disable
"Do not use AUTO_INCREMENT for zero values"

Make sure to open the file called "config.php" it is in the "inc" folder
change the values to your SQL Settings and Database Name

```sh
$ sudo chmod 777 config.json
```
then visit the login.php page. and login. Default Login Details are
```
Username: Master
Password: Shadi1999
```

## Pictures
<img src="https://i.imgur.com/fmjDpVR.png">
<img src="https://i.imgur.com/IYTtmwj.png">
<img src="https://i.imgur.com/fdgIL1k.png">
<img src="https://i.imgur.com/dWqPwU0.png">

### Features
* Auto Detects Google Home Devices
* Customizable Athan Times and Location 
* Change each salat volume
* Change Athan Audio
* 15 Minute Salat warning
* Secure Login
* Disable and Enable Diffrent Salats
* Recites Surah Kahaf after Duhur athan on Fridays

### Todos

 - Unknown, give suggestions

## READ
DO NOT SELL THIS IN ANY FORM, SERVICE, PRODUCT OR ANYTHING
DO NOT USE FOR COMMERCIAL USSAGE
