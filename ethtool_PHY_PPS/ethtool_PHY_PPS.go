package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/safchain/ethtool"
)


func GetStats(e  *ethtool.Ethtool, name string) (uint64, uint64, uint64, uint64, uint64) {
	stats, err := e.Stats(name)
        if err != nil {
                panic(err.Error())
        }
	return  stats["rx_bytes"], stats["rx_packets"], stats["tx_bytes"], stats["tx_packets"], stats["rx_out_of_buffer"]


}


type OneIntf struct {
	e  *ethtool.Ethtool
	intf		string
	display_rx_bytes uint64
	display_rx_packets uint64
	display_tx_bytes uint64
	display_tx_packets uint64
	display_rx_drops uint64
	rx_bytes uint64
	rx_packets uint64
	tx_bytes uint64
	tx_packets uint64
	rx_drops uint64
}

func (self *OneIntf) do() {

	stats, err := self.e.Stats(self.intf)
        if err != nil {
                panic(err.Error())
        }
	self.display_rx_bytes   = stats["rx_bytes"] - self.rx_bytes
	self.display_rx_packets = stats["rx_packets"] - self.rx_packets
	self.display_tx_bytes   = stats["tx_bytes"] - self.tx_bytes
	self.display_tx_packets = stats["tx_packets"] - self.tx_packets
	self.display_rx_drops   = stats["rx_out_of_buffer"] - self.rx_drops
	self.rx_bytes   = stats["rx_bytes"]
	self.rx_packets = stats["rx_packets"]
	self.tx_bytes   = stats["tx_bytes"]
	self.tx_packets = stats["tx_packets"]
	self.rx_drops   = stats["rx_out_of_buffer"]

	//fmt.Printf("\033[0;0H")
	//fmt.Println(self.display_rx_bytes, self.display_rx_packets, self.display_tx_bytes, self.display_tx_packets, self.display_rx_drops)


}


func (self *OneIntf) mainLoop() {
	//init
	stats, err := self.e.Stats(self.intf)
	self.rx_bytes   = stats["rx_bytes"]
	self.rx_packets = stats["rx_packets"]
	self.tx_bytes   = stats["tx_bytes"]
	self.tx_packets = stats["tx_packets"]
	self.rx_drops   = stats["rx_out_of_buffer"]

        if err != nil {
                panic(err.Error())
        }

	for {
		go self.do()
		time.Sleep(1000 * time.Millisecond)
	}

}

func main() {

	print("\033[H\033[2J")



	if len(os.Args) < 2 {
		log.Fatal("Please specify at least one interface")
	}


	e, err := ethtool.NewEthtool()
	if err != nil {
		panic(err.Error())
	}
	defer e.Close()

	baseLine := 4

	pos := fmt.Sprintf("\033[%d;14H", baseLine)
	//fmt.Printf("\033[4;14H")
	fmt.Printf(pos)
	fmt.Println("RX              TX                  ERRORS")


	intfs := make([]OneIntf, 0)
	numIntfs := len(os.Args)-1

	for i:=0; i < numIntfs;  i++ {

		intf := OneIntf{e: e, intf: os.Args[i+1]}
		intfs = append(intfs, intf)
		go intf.mainLoop()
	}

	for {
		time.Sleep(1000 * time.Millisecond)
		line := baseLine+1;
		pos = fmt.Sprintf("\033[%d;14H", line)
		fmt.Printf(pos)
		for i := 0 ; i < numIntfs; i++ {
			//fmt.Println(intfs[i].display_rx_bytes, intfs[i].display_rx_packets, intfs[i].display_tx_bytes, intfs[i].display_tx_packets, intfs[i].display_rx_drops)
		}
	}


	select  {}
	//GetStats(e, intf)

}
/*
func main() {
	name := flag.String("interface", "", "Interface name")
	flag.Parse()

	if *name == "" {
		log.Fatal("interface is not specified")
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		panic(err.Error())
	}
	defer e.Close()

	features, err := e.Features(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("features: %+v\n", features)

	stats, err := e.Stats(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("stats: %+v\n", stats)

	busInfo, err := e.BusInfo(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("bus info: %+v\n", busInfo)

	drvr, err := e.DriverName(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("driver name: %+v\n", drvr)

	cmdGet, err := e.CmdGetMapped(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("cmd get: %+v\n", cmdGet)

	msgLvlGet, err := e.MsglvlGet(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("msg lvl get: %+v\n", msgLvlGet)

	drvInfo, err := e.DriverInfo(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("drvrinfo: %+v\n", drvInfo)

	permAddr, err := e.PermAddr(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("permaddr: %+v\n", permAddr)

	eeprom, err := e.ModuleEepromHex(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("module eeprom: %+v\n", eeprom)
}
*/
