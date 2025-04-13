package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type OsmAddress struct {
	Id     string
	Street string
	Build  string
	Levels int
	LatLng dto.LatLng
}

type OsmWayTag struct {
	K string `xml:"k,attr"`
	V string `xml:"v,attr"`
}

type OsmWay struct {
	Id   string      `xml:"id,attr"`
	Tags []OsmWayTag `xml:"tag"`
}

type OsmData struct {
	XMLName xml.Name `xml:"osm"`
	Ways    []OsmWay `xml:"way"`
}

type OsmLatLngResult struct {
	Lat float64 `xml:"lat,attr"`
	Lng float64 `xml:"lon,attr"`
}

type OsmLatLngData struct {
	XMLName xml.Name        `xml:"reversegeocode"`
	Result  OsmLatLngResult `xml:"result"`
}

func (ths *Simulator) getOsmLatLng(id string) (dto.LatLng, error) {
	request := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=xml&osm_type=W&osm_id=%s&zoom=18&addressdetails=1", id)
	resp, err := http.Get(request)
	if err != nil {
		return dto.LatLng{}, fmt.Errorf("can't get latlng: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dto.LatLng{}, fmt.Errorf("can't read response: %w", err)
	}

	var osmLatLng OsmLatLngData

	err = xml.Unmarshal(data, &osmLatLng)
	if err != nil {
		return dto.LatLng{}, fmt.Errorf("invalid response: %w", err)
	}

	return dto.LatLng{
		Lat: osmLatLng.Result.Lat,
		Lng: osmLatLng.Result.Lng,
	}, nil
}

func (ths *Simulator) getOsmData(leftUpper, rightBottom dto.LatLng, minLevel, max int) ([]OsmAddress, error) {
	ths.log.Info("Downloading OSM data...")

	request := fmt.Sprintf("https://api.openstreetmap.org/api/0.6/map?bbox=%f,%f,%f,%f",
		leftUpper.Lat,
		leftUpper.Lng,
		rightBottom.Lat,
		rightBottom.Lng,
	)
	ths.log.Info(request)

	resp, err := http.Get(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var osmData OsmData
	err = xml.Unmarshal(data, &osmData)
	if err != nil {
		return nil, err
	}

	var addresses []OsmAddress

	for _, w := range osmData.Ways {

		addr := OsmAddress{}

		for _, tag := range w.Tags {
			switch tag.K {
			case "addr:housenumber":
				addr.Build = tag.V
			case "addr:street":
				addr.Street = tag.V
			case "building:levels":
				levels, err := strconv.Atoi(tag.V)
				if err != nil {
					ths.log.Error("Invalid level: %s", tag.V)
					continue
				}
				addr.Levels = levels
			}
		}
		if addr.Levels > 0 && len(addr.Street) > 0 && len(addr.Build) > 0 && addr.Levels >= minLevel {
			addr.Id = w.Id
			addresses = append(addresses, addr)
			if max > 0 && len(addresses) >= max {
				break
			}
		}
	}

	ths.log.Info("Getting LatLng for %d addresses...", len(addresses))

	addressCache := ths.getAddressCache()

	for i, a := range addresses {
		addr, err := ths.findInAddressCache(addressCache, a.Id)
		if err != nil {
			latlng, err := ths.getOsmLatLng(a.Id)
			if err != nil {
				ths.log.Error("Can't get LatLng for %s: %v", a.Id, err)
				continue
			}

			ths.log.Info("%d/%d (direct)", i+1, len(addresses))
			addresses[i].LatLng = latlng
			addressCache = append(addressCache, addresses[i])
		} else {
			ths.log.Info("%d/%d (cache)", i+1, len(addresses))
			addresses[i].LatLng = addr.LatLng
		}
	}

	err = ths.saveAddressCache(addressCache)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (ths *Simulator) findInAddressCache(addresses []OsmAddress, id string) (OsmAddress, error) {
	for _, a := range addresses {
		if a.Id == id {
			return a, nil
		}
	}
	return OsmAddress{}, fmt.Errorf("not found in cache")
}

func (ths *Simulator) saveAddressCache(addresses []OsmAddress) error {
	var data []byte
	data, err := json.Marshal(addresses)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("addr_cache.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ths *Simulator) getAddressCache() []OsmAddress {
	data, err := os.ReadFile("addr_cache.json")
	if err != nil {
		return []OsmAddress{}
	}

	var addresses []OsmAddress
	err = json.Unmarshal(data, &addresses)
	if err != nil {
		return []OsmAddress{}
	}

	return addresses
}
