package main

/*
	- nohup ./Shadi &
	- chmod -R 777 ~/

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/evalphobia/google-home-client-go/googlehome"
)

//Athan grabs main data header from Json
type Athan struct {
	Data YtData `json:"data"`
}

//YtData grabs secondary header under "Athan"
type YtData struct {
	Timings YtTime `json:"timings"`
}

//YtTime Gets 3rd Header from Json file, which is where the athan times are located
type YtTime struct {
	F string `json:"Fajr"`
	D string `json:"Dhuhr"`
	A string `json:"Asr"`
	M string `json:"Maghrib"`
	I string `json:"Isha"`
}

//Config Get Config settings from config.json file
type Config struct {
	Settings struct {
		Name     string `json:"Name"`
		Language string `json:"Language"`
		Accent   string `json:"Accent"`
		Athan    string `json:"Athan"`
	}

	Connection struct {
		IP   string `json:"IP"`
		Port int    `json:"Port"`
	}

	Prayers struct {
		Fajir  bool `json:"Fajir"`
		Duhur  bool `json:"Duhur"`
		Asr    bool `json:"Asr"`
		Magrib bool `json:"Magrib"`
		Isha   bool `json:"Isha"`
	}

	Audio struct {
		Athan  string `json:"Athan"`
		Recite string `json:"Recite"`
	}

	Location struct {
		City     string `json:"City"`
		Country  string `json:"Country"`
		State    string `json:"State"`
		TimeZone string `json:"TimeZone"`
	}

	Calculation struct {
		Method int `json:"Method"`
	}

	Volume struct {
		Connection bool    `json:"Connection"`
		Default    float64 `json:"Default"`
		Fajir      float64 `json:"Fajir"`
		Duhur      float64 `json:"Duhur"`
		Asr        float64 `json:"Asr"`
		Magrib     float64 `json:"Magrib"`
		Isha       float64 `json:"Isha"`
	}

	Options struct {
		Whisper bool `json:"Whisper"`
		Recite  bool `json:"Recite"`
	}
}

//Y Gets assigned from Athan
//var Y Athan

//Split API
const (
	MainAPI    string = "http://api.aladhan.com/v1/timingsByCity?city="
	CountryAPI string = "&country="
	StateAPI   string = "&state="
	MethodAPI  string = "&method="
)

var Meth = strconv.Itoa(config.Calculation.Method)

var config Config

func main() {

	var err error
	//Connect to Json file for settings and paramaters
	config, err = LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error importing config.json file", err)
	}

	//Connect to Google Home
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: config.Connection.IP,
		Lang:     config.Settings.Language,
		Accent:   config.Settings.Accent,
		Port:     42038,
	})

	if err != nil {
		panic(err)
	} else {
		// Sets to device to default volume
		cli.SetVolume(config.Volume.Default)

		//Echos to device to tell if users its Connected
		if config.Volume.Connection == true {
			cli.Notify("Successfully Connected.")
			//cli.Play(config.Audio.Athan)
		}
		ConnectedTo()
	}

	//Call Athan API
	Y := ACal()

	for range time.Tick(time.Second * 25) {
		//Grab Updated Config Files
		config, _ := LoadConfig("config.json")

		//Get Local time test
		t := time.Now()
		location, err := time.LoadLocation(config.Location.TimeZone)
		if err != nil {
			fmt.Println(err)
		}

		CurrentTime := fmt.Sprint(t.In(location).Format("15:04"))

		//Check if friday
		day := time.Now().Weekday()
		CurrentDay := fmt.Sprint(day)

		//Checks if its time for Fajir
		if Y.Data.Timings.F == CurrentTime {
			if config.Prayers.Fajir == true {
				//fmt.Println("Time for Fajir")
				cli.SetVolume(config.Volume.Fajir)
				cli.Play(config.Audio.Athan)
				time.Sleep(4 * time.Minute)
			}
		}

		//Checks if its time for Duhur
		if Y.Data.Timings.D == CurrentTime {
			if config.Prayers.Duhur == true {
				//fmt.Println("Time for Duhur")
				cli.SetVolume(config.Volume.Duhur)
				cli.Play(config.Audio.Athan)
				time.Sleep(4 * time.Minute)
			}
		}

		if config.Prayers.Duhur == true {
			//fmt.Println("Time for Duhur")
			cli.SetVolume(config.Volume.Duhur)
			cli.Play(config.Audio.Athan)
			time.Sleep(4 * time.Minute)
		}

		//Checks if the day is Friday
		if config.Options.Recite == true {
			if CurrentDay == "Friday" {
				//cli.Notify("I will begin reciting Quran.")
				time.Sleep(5 * time.Second)
				cli.Play(config.Audio.Recite)
				time.Sleep(30 * time.Minute)
			}
		}

		//Checks if its time for Asr
		if Y.Data.Timings.A == CurrentTime {
			if config.Prayers.Asr == true {
				//fmt.Println("Time for Asr")
				cli.SetVolume(config.Volume.Asr)
				cli.Play(config.Audio.Athan)
				time.Sleep(4 * time.Minute)
			}
		}

		//Checks if its time for Magrib
		if Y.Data.Timings.M == CurrentTime {
			if config.Prayers.Magrib == true {
				//fmt.Println("Time for Magrib")
				cli.SetVolume(config.Volume.Magrib)
				cli.Play(config.Audio.Athan)
				time.Sleep(4 * time.Minute)

			}
		}

		//Checks if time for Isha
		if Y.Data.Timings.I == CurrentTime {
			if config.Prayers.Isha == true {
				//fmt.Println("Time for Isha")
				cli.SetVolume(config.Volume.Isha)
				cli.Play(config.Audio.Athan)
				time.Sleep(4 * time.Minute)
				ACal() //Recall Json Data
			}

		}
	} // End Loop

}

//ACal API Function
func ACal() Athan {
	var AthanAPI = MainAPI + config.Location.City + CountryAPI + config.Location.Country + StateAPI + config.Location.State + MethodAPI + Meth
	FormatAPI := fmt.Sprintf(AthanAPI)

	//fmt.Println(AthanAPI)

	resp, err := http.Get(FormatAPI)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	var Y Athan
	err = json.Unmarshal(body, &Y)
	if err != nil {
		log.Fatal(err)
	}

	return Y
}

//LoadConfig file
func LoadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)

	defer configFile.Close()
	if err != nil {
		return config, err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

//ConnectedTo gives information of connected google home and its basic paramaters
func ConnectedTo() {
	config, _ := LoadConfig("config.json")

	fmt.Println("Device Connected:", time.Now())
	fmt.Println("Connected to Device:", config.Settings.Name)
	fmt.Println("IP Address:", config.Connection.IP)
	fmt.Println("Using Lanuage:", config.Settings.Language)
	fmt.Println("Using Accent:", config.Settings.Accent)
	fmt.Println("Default Volume Set at", config.Volume.Default)

	//Calculation Method
	MethodV()

}

//MethodV Find out what Calculation method is being used
func MethodV() {
	config, _ := LoadConfig("config.json")

	switch config.Calculation.Method {
	case 0:
		//fmt.Println("Using Calculation Method: Shia Ithna-Ansari")
	case 1:
		fmt.Println("Using Calculation Method: University of Islamic Sciences, Karachi")
	case 2:
		fmt.Println("Using Calculation Method: Islamic Society of North America")
	case 3:
		fmt.Println("Using Calculation Method: Muslim World League")
	case 4:
		fmt.Println("Using Calculation Method: Umm Al-Qura University, Makkah")
	case 5:
		fmt.Println("Using Calculation Method: Egyptian General Authority of Survey")
	case 7:
		fmt.Println("Using Calculation Method: Institute of Geophysics, University of Tehran")
	case 8:
		fmt.Println("Using Calculation Method: Gulf Region")
	case 9:
		fmt.Println("Using Calculation Method: Kuwait")
	case 10:
		fmt.Println("Using Calculation Method: Qatar")
	case 11:
		fmt.Println("Using Calculation Method: Majlis Ugama Islam Singapura, Singapore")
	case 12:
		fmt.Println("Using Calculation Method: Union Organization islamic de France")
	case 13:
		fmt.Println("Using Calculation Method: Diyanet İşleri Başkanlığı, Turkey")
	default:
		fmt.Println("Other option choosen")
	}
}
