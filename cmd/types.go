package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type IPInfo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ADSBOneResponse struct {
	Ac    []Ac   `json:"ac"`
	Msg   string `json:"msg"`
	Now   int64  `json:"now"`
	Total int    `json:"total"`
	Ctime int64  `json:"ctime"`
	Ptime int    `json:"ptime"`
}
type Ac struct {
	Hex            string      `json:"hex"`
	Type           string      `json:"type"`
	Flight         string      `json:"flight,omitempty"`
	R              string      `json:"r,omitempty"`
	T              string      `json:"t,omitempty"`
	Desc           string      `json:"desc,omitempty"`
	OwnOp          string      `json:"ownOp,omitempty"`
	AltBaro        StringOrInt `json:"alt_baro"` // not reliably an int
	AltGeom        int         `json:"alt_geom,omitempty"`
	Gs             float64     `json:"gs,omitempty"`
	Track          float64     `json:"track,omitempty"`
	BaroRate       int         `json:"baro_rate,omitempty"`
	Squawk         string      `json:"squawk,omitempty"`
	Emergency      string      `json:"emergency,omitempty"`
	Category       string      `json:"category,omitempty"`
	NavQnh         float64     `json:"nav_qnh,omitempty"`
	NavAltitudeMcp int         `json:"nav_altitude_mcp,omitempty"`
	NavHeading     float64     `json:"nav_heading,omitempty"`
	Lat            float64     `json:"lat"`
	Lon            float64     `json:"lon"`
	Nic            int         `json:"nic"`
	Rc             int         `json:"rc"`
	SeenPos        float64     `json:"seen_pos"`
	Version        int         `json:"version,omitempty"`
	NicBaro        int         `json:"nic_baro,omitempty"`
	NacP           int         `json:"nac_p,omitempty"`
	NacV           int         `json:"nac_v,omitempty"`
	Sil            int         `json:"sil,omitempty"`
	SilType        string      `json:"sil_type,omitempty"`
	Gva            int         `json:"gva,omitempty"`
	Sda            int         `json:"sda,omitempty"`
	Alert          int         `json:"alert,omitempty"`
	Spi            int         `json:"spi,omitempty"`
	Mlat           []any       `json:"mlat"`
	Tisb           []any       `json:"tisb"`
	Messages       int         `json:"messages"`
	Seen           float64     `json:"seen"`
	Rssi           float64     `json:"rssi"`
	Dst            float64     `json:"dst"`
	Dir            float64     `json:"dir"`
	Year           string      `json:"year,omitempty"`
	GeomRate       int         `json:"geom_rate,omitempty"`
	NavModes       []string    `json:"nav_modes,omitempty"`
	DbFlags        int         `json:"dbFlags,omitempty"`
	CalcTrack      int         `json:"calc_track,omitempty"`
	TrueHeading    float64     `json:"true_heading,omitempty"`
	NavAltitudeFms int         `json:"nav_altitude_fms,omitempty"`
	MagHeading     float64     `json:"mag_heading,omitempty"`
}

func (a *Ac) getAltitude() int {
	return a.AltBaro.Value
}
func (a *Ac) getFlight() string {
	return strings.ToUpper(strings.TrimSpace(a.Flight))
}

type StringOrInt struct {
	Value int
}

func (si *StringOrInt) UnmarshalJSON(data []byte) error {
	var str string
	var num int
	if err := json.Unmarshal(data, &str); err == nil {
		// Force to 0 when not a valid int
		si.Value = 0
		return nil
	} else if err := json.Unmarshal(data, &num); err == nil {
		si.Value = num
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into StringOrInt", data)
}
