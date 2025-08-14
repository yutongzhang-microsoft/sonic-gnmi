package show_client

import (
    "sort"
    "strings"
    "strconv"
    "github.com/facette/natsort"
	log "github.com/golang/glog"
	sdc "github.com/sonic-net/sonic-gnmi/sonic_data_client"
)

func isRoleInternal(role string) bool {
    return role == "Int" || role == "Inb" || role == "Rec" || role == "Dpc"
}

func isFrontPanelPort(iface string, role string) bool {
    if !strings.HasPrefix(iface, "Ethernet") {
        return false
    }
    if strings.HasPrefix(iface, "Ethernet-BP") || strings.HasPrefix(iface, "Ethernet-IB") || strings.HasPrefix(iface, "Ethernet-Rec") {
        return false
    }
    if strings.Contains(iface, ".") {
        return false
    }
    return !isRoleInternal(role)
}

func isValidPhysicalPort(iface string) (bool, error) {
    queries := [][]string{
	    {"APPL_DB", "PORT_TABLE"},
	}
	portTable, err := GetMapFromQueries(queries)
	if err != nil {
		log.Errorf("Unable to pull data for queries %v, got err %v", queries, err)
		return false, err
	}
    role := GetFieldValueString(portTable, iface, defaultMissingCounterValue, "role")
    return isFrontPanelPort(iface, role), nil
}

func getLogicalToPhysical(logicalPort string) (string, error) {
    logicalToPhysical := make(map[string][]int)
    logical := []string{}

    queries := [][]string{
        {"CONFIG_DB", "PORT"},
    }
    portTable, err := GetMapFromQueries(queries)
	if err != nil {
		log.Errorf("Unable to pull data for queries %v, got err %v", queries, err)
		return nil, err
	}
	for key := range portTable {
        parts := strings.SplitN(key, ":", 2)

	    var iface string
	    if len(parts) == 2 {
            iface = strings.TrimSpace(parts[1])
	    } else {
            iface = strings.TrimSpace(parts[0])
	    }

	    if isFrontPanel(iface, GetFieldValueString(portTable, iface, defaultMissingCounterValue, "role")){
	        logical = append(logical, iface)
	    }
	}

	for _, intfName := range logical {
        if _, ok := portTable[intfName]["index"]; ok{
            fpPortIndex := strconv.Atoi(portTable[intfName]["index"])
            logicalToPhysical[intfName] = []int{fpPortIndex}
        }
	}

    return logicalToPhysical[logicalPort]
}

func getFirstSubPort(logicalPort string) {
    physicalPort := getLogicalToPhysical(logicalPort)
}

func convertInterfaceSfpInfoToCliOutputString(iface string) {
    firstPort := getFirstSubPort(iface)

}

func getEEPROM(options sdc.OptionMap) (map[string]string, error){
    var intf string
	if v, ok := options["interface"].Strings(); ok {
		intf = v
	}

    var queries [][]string
	if intf == "" {
		queries = [][]string{
			{"APPL_DB", "PORT_TABLE"},
		}
	} else {
// 		queries = [][]string{
// 			{"STATE_DB", "TRANSCEIVER_STATUS_SW", intf},
// 		}
	}

// 	portTableKeys  := []string{}

	portTable, err := GetMapFromQueries(queries)
	if err != nil {
		log.Errorf("Unable to pull data for queries %v, got err %v", queries, err)
		return nil, err
	}

	intfEEPROM := make(map[string]string)
	for key := range portTable {
	    parts := strings.SplitN(key, ":", 2)

	    var iface string
	    if len(parts) == 2 {
            iface = strings.TrimSpace(parts[1])
	    } else {
            iface = strings.TrimSpace(parts[0])
	    }

        if iface != "" {
            ok, err := isValidPhysicalPort(iface)
            if err != nil {
                return nil, err
            }
            if ok {
                intfEEPROM[iface] = convertInterfaceSfpInfoToCliOutputString(iface)
            }
        }
	}
	return intfEEPROM, nil
}

