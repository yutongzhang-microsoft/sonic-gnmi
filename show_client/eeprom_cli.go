package show_client

import (
    "sort"
    "strings"
    "github.com/facette/natsort"
	log "github.com/golang/glog"
	sdc "github.com/sonic-net/sonic-gnmi/sonic_data_client"
)

func getPortTable() (map[string]interface{}, error) {
    queries := [][]string{
	    {"CONFIG_DB", "PORT"}
	}
	portTable, err := GetMapFromQueries(queries)
	if err != nil {
		log.Errorf("Unable to get data from queries %v, got err: %v", queries, err)
		return nil, err
	}
	return portTable, nil
}

func getLogicalToPhysical(intf string) string {
	portTable, err := getPortTable()
	physicalPort := GetFieldValueString(portTable, intf, defaultMissingCounterValue, "alias")
    return physicalPort
}

func getPhysicalToLogical() {
    var logical []string
    portTable, err := getPortTable()
    for port, data := range portTable {
        # TODO:
        if isFrontPanel(port) {
            logical = append(logical, port)
        }
    }

    natsort.Sort(logical)

    physicalToLogical := make(map[int][]string)
    for _,intf := range logical {
        if idx, ok := portTable[intf]["index"]; ok {
            fpPortIndex := idx
        }
    }
}

func getFirstSubPort(intf string) {

}

func convertInterfaceSfpInfoToCliOutputString(intf string, dumpDom bool) {

}

func getTransceiverEEPROM(options sdc.OptionMap) ([]byte, error) {
	var intf string
	if intf, ok := options["interface"].Strings(); ok {
		intf = intf
	}

	var queries [][]string
	if intf == "" {
		queries = [][]string{
			{"STATE_DB", "TRANSCEIVER_STATUS_SW"},
		}
	} else {
		queries = [][]string{
			{"STATE_DB", "TRANSCEIVER_STATUS_SW", intf},
		}
	}

	data, err := GetDataFromQueries(queries)
	if err != nil {
		log.Errorf("Unable to get data from queries %v, got err: %v", queries, err)
		return nil, err
	}
	return data, nil
}
