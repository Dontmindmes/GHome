package main

/*
	- nohup ./Shadi &
	- chmod -R 777 ~/

	POSSIBLY ADD TUNE FUNCTION
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evalphobia/google-home-client-go/googlehome"
	"github.com/gojektech/heimdall/httpclient"
	"github.com/micro/mdns"
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

	Prayers struct {
		Fajir  bool `json:"Fajir"`
		Duhur  bool `json:"Duhur"`
		Asr    bool `json:"Asr"`
		Magrib bool `json:"Magrib"`
		Isha   bool `json:"Isha"`
	}

	Audio struct {
		Athan  string `json:"Athan"`
		Fajir  string `json:"Fajir"`
		Recite string `json:"Recite"`
	}

	Location struct {
		City     string `json:"City"`
		Country  string `json:"Country"`
		State    string `json:"State"`
		TimeZone string `json:"TimeZone"`
	}

	Calculation struct {
		Method int    `json:"Method"`
		School string `json:"School"`
	}

	Volume struct {
		Default float64 `json:"Default"`
		Fajir   float64 `json:"Fajir"`
		Duhur   float64 `json:"Duhur"`
		Asr     float64 `json:"Asr"`
		Magrib  float64 `json:"Magrib"`
		Isha    float64 `json:"Isha"`
	}

	Options struct {
		Recite     bool `json:"Recite"`
		Alert      bool `json:"Alert"`
		Connection bool `json:"Connection"`
	}
}

const (
	googleCastServiceName = "_googlecast._tcp"
	googleHomeModelInfo   = "md=Google Home"
)

type GoogleHomeInfo struct {
	IP   string
	Port int
}

var (
	CPort int
	CIP   string
)

//Split API
const (
	MainAPI    = "http://api.aladhan.com/v1/timingsByCity?city="
	CountryAPI = "&country="
	StateAPI   = "&state="
	MethodAPI  = "&method="
	SchoolAPI  = "&school="
)

var config Config
var cli *googlehome.Config

func main() {
	LookupHomeIP()
	var err error

	//Connect to Json file for settings and paramaters
	config, err = LoadConfig("config.json")
	if err != nil {
		CheckOFF()
		log.Fatal("Error importing config.json file", err)
	}

	//Connect to Google Home
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: CIP,
		Lang:     config.Settings.Language,
		Accent:   config.Settings.Accent,
		Port:     CPort,
	})

	if err != nil {
		CheckOFF()
		panic(err)
	}

	CheckOn()
	// Sets to device to default volume
	cli.SetVolume(config.Volume.Default)
	//Echos to device to tell if users its Connected
	if config.Options.Connection == true {
		cli.Notify("Successfully Connected.")
	}
	ConnectedTo()

	//Call Athan API
	Y := ACal()

	for range time.Tick(time.Second * 10) {
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
		day := t.In(location).Weekday()
		CurrentDay := fmt.Sprint(day)

		//Duhur
		pd, _ := time.Parse("15:04", Y.Data.Timings.D)
		pd = pd.Add(time.Minute * time.Duration(-15))
		pad := fmt.Sprintf(pd.Format("15:04"))

		//Asr
		pa, _ := time.Parse("15:04", Y.Data.Timings.A)
		pa = pa.Add(time.Minute * time.Duration(-15))
		pas := fmt.Sprintf(pa.Format("15:04"))

		//Magrib
		pm, _ := time.Parse("15:04", Y.Data.Timings.M)
		pm = pm.Add(time.Minute * time.Duration(-15))
		pam := fmt.Sprintf(pm.Format("15:04"))

		//Isha
		pi, _ := time.Parse("15:04", Y.Data.Timings.I)
		pi = pi.Add(time.Minute * time.Duration(-15))
		pai := fmt.Sprintf(pi.Format("15:04"))

		//Checks if its time for Fajir
		if Y.Data.Timings.F == CurrentTime && config.Prayers.Fajir {
			cli.SetVolume(config.Volume.Fajir)
			cli.Play(config.Audio.Fajir)
			time.Sleep(4 * time.Minute)
		}

		//PreAlert for Duhur
		if config.Prayers.Duhur && config.Options.Alert && pad == CurrentTime {
			cli.SetVolume(config.Volume.Default)
			cli.Notify(" ")
			time.Sleep(15 * time.Second)
		}

		//Checks if its time for Duhur
		if Y.Data.Timings.D == CurrentTime && config.Prayers.Duhur {
			cli.SetVolume(config.Volume.Duhur)
			cli.Play(config.Audio.Athan)
			time.Sleep(6 * time.Minute)

			//Checks if the day is Friday
			if config.Options.Recite && CurrentDay == "Friday" {
				//cli.Notify("I will begin reciting Quran.")
				time.Sleep(5 * time.Second)
				cli.Play(config.Audio.Recite)
				time.Sleep(30 * time.Minute)
			}
		}

		//PreAlert for Asr
		if config.Prayers.Asr && config.Options.Alert && pas == CurrentTime {
			cli.SetVolume(config.Volume.Default)
			cli.Notify(" ")
			time.Sleep(15 * time.Second)
		}

		//Checks if its time for Asr
		if Y.Data.Timings.A == CurrentTime && config.Prayers.Asr {
			cli.SetVolume(config.Volume.Asr)
			cli.Play(config.Audio.Athan)
			time.Sleep(4 * time.Minute)
		}

		//PreAlert for Magrib
		if config.Prayers.Magrib && config.Options.Alert && pam == CurrentTime {
			cli.SetVolume(config.Volume.Default)
			cli.Notify(" ")
			time.Sleep(15 * time.Second)
		}

		//Checks if its time for Magrib
		if Y.Data.Timings.M == CurrentTime && config.Prayers.Magrib {
			cli.SetVolume(config.Volume.Magrib)
			cli.Play(config.Audio.Athan)
			time.Sleep(4 * time.Minute)
		}

		//PreAlert for Isha
		if config.Prayers.Isha && config.Options.Alert && pai == CurrentTime {
			cli.SetVolume(config.Volume.Default)
			cli.Notify(" ")
			time.Sleep(15 * time.Second)
		}

		//Checks if time for Isha
		if Y.Data.Timings.I == CurrentTime && config.Prayers.Isha {
			cli.SetVolume(config.Volume.Isha)
			cli.Play(config.Audio.Athan)
			time.Sleep(4 * time.Minute)
			ACal() //Recall Json Data
		}
	} // End Loop

}

//LookupHomeIP Scans internet for Google Cast System
func LookupHomeIP() []*GoogleHomeInfo {
	
	entriesCh := make(chan *mdns.ServiceEntry, 2)
	results := []*GoogleHomeInfo{}

	go func() {
		for entry := range entriesCh {
			if strings.Contains(entry.Name, "Cast") {
				fmt.Printf("[INFO] ServiceEntry Group Cast detected: [%s:%d] \n", entry.AddrV4, entry.Port)
				CPort = entry.Port
				CIP = fmt.Sprintf(entry.AddrV4.String())
			}
			for _, field := range entry.InfoFields {
				if strings.HasPrefix(field, googleHomeModelInfo) {
					results = append(results, &GoogleHomeInfo{entry.AddrV4.String(), entry.Port})
				}
			}
		}
	}()

	mdns.Lookup(googleCastServiceName, entriesCh)
	close(entriesCh)

	return results
}

//ACal API Function
func ACal() Athan {
	var MethConv = strconv.Itoa(config.Calculation.Method)
	var AthanAPI = MainAPI + config.Location.City + CountryAPI + config.Location.Country + StateAPI + config.Location.State + MethodAPI + MethConv + SchoolAPI + config.Calculation.School

	//var AthanAPI = MainAPI + config.Location.City + CountryAPI + config.Location.Country + StateAPI + config.Location.State + MethodAPI + Meth + SchoolAPI + config.Calculation.School
	FormatAPI := fmt.Sprintf(AthanAPI)

	// Create a new HTTP client with a default timeout
	timeout := 3000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	// Use the clients GET method to create and execute the request
	resp, err := client.Get(FormatAPI, nil)
	if err != nil {
		defer CheckOFF()
		panic(err)
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
	fmt.Println("IP Address:", CIP)
	fmt.Println("Port:", CPort)
	fmt.Println("Using Lanuage:", config.Settings.Language)
	fmt.Println("Using Accent:", config.Settings.Accent)
	fmt.Println("Default Volume Set at", config.Volume.Default)

	//Calculation Method
	CalcMethodPrint()
}

//CalcMethodPrint Find out what Calculation method is being used
func CalcMethodPrint() {
	config, _ := LoadConfig("config.json")

	var Using string = "Using Calculation Method: "

	switch config.Calculation.Method {
	case 0:
		//fmt.Println("Using Calculation Method: Shia Ithna-Ansari")
	case 1:
		fmt.Println(Using + " University of Islamic Sciences, Karachi")
	case 2:
		fmt.Println(Using + " Islamic Society of North America")
	case 3:
		fmt.Println(Using + " Muslim World League")
	case 4:
		fmt.Println(Using + " Umm Al-Qura University, Makkah")
	case 5:
		fmt.Println(Using + " Egyptian General Authority of Survey")
	case 7:
		fmt.Println(Using + " Institute of Geophysics, University of Tehran")
	case 8:
		fmt.Println(Using + " Gulf Region")
	case 9:
		fmt.Println(Using + " Kuwait")
	case 10:
		fmt.Println(Using + " Qatar")
	case 11:
		fmt.Println(Using + " Majlis Ugama Islam Singapura, Singapore")
	case 12:
		fmt.Println(Using + " Union Organization islamic de France")
	case 13:
		fmt.Println(Using + " Diyanet İşleri Başkanlığı, Turkey")
	case 14:
		fmt.Println(Using + " Spiritual Administration of Muslims of Russia")
	default:
		fmt.Println("Other option choosen")
	}
}

func CheckOn() {
	var err = os.Remove("Status.txt")
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create("Status.txt")
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.WriteString("ON")

	defer file.Close()
}

func CheckOFF() {
	file, _ := os.Create("Status.txt")

	s, err := file.WriteString("OFF")
	if err != nil {
		fmt.Println(s, err)
	}
	defer file.Close()
}
