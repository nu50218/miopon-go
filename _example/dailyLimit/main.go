package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/nu50218/miopon-go/coupon"

	"github.com/nu50218/miopon-go"
)

type Settings struct {
	DeveloperID    string `json:"developer_id"`
	AccessToken    string `json:"access_tokens"`
	IntervalSecond int    `json:"interval_second"`
	VolumeLimit    int    `json:"volume_limit"`
}

const settingsFilename = "settings.json"

func main() {
	var isInit bool
	flag.BoolVar(&isInit, "init", false, "Make settings.json in the working directory")
	var developerID string
	flag.StringVar(&developerID, "id", "", "developerID")
	flag.Parse()

	if isInit {
		if developerID == "" {
			fmt.Print("developerID:")
			fmt.Scan(&developerID)
		}

		var redirectURI, state string
		fmt.Print("redirectURI:")
		fmt.Scan(&redirectURI)
		fmt.Print("state:")
		fmt.Scan(&state)

		settings := &Settings{
			DeveloperID:    developerID,
			AccessToken:    "",
			IntervalSecond: 60,
		}
		b, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile(settingsFilename, b, 0644); err != nil {
			log.Fatalln(err)
		}

		client := miopon.New(developerID, &miopon.Settings{})
		fmt.Println("Created settings file. Get access tokens from the following URL and add them to the settings file.")
		fmt.Println(client.MakeAuthorizationURL(redirectURI, state))
		return
	}

	b, err := ioutil.ReadFile(settingsFilename)
	if err != nil {
		log.Fatalln(err)
	}

	settings := &Settings{}
	if err := json.Unmarshal(b, settings); err != nil {
		log.Fatalln(err)
	}

	client := miopon.New(settings.DeveloperID, &miopon.Settings{
		GetCouponInterval: 12 * time.Second,
		PutCouponInterval: 60 * time.Second,
		GetPacketInterval: 12 * time.Second,
	})

	for {
		check(settings.AccessToken, client, settings.VolumeLimit)
		time.Sleep(time.Duration(settings.IntervalSecond) * time.Second)
	}

}

func check(accessToken string, client *miopon.Client, volumeLimit int) {
	body, statusCode, err := client.GetPacket(accessToken)
	if err != nil {
		log.Println(err)
		return
	}
	if statusCode != http.StatusOK {
		log.Println(statusCode, body.ReturnCode)
		return
	}

	couponInfo := []*coupon.CouponInfo{}
	for _, packetLogInfo := range body.PacketLogInfo {
		c := coupon.CouponInfo{}
		for _, hdoInfo := range packetLogInfo.HdoInfo {
			if len(hdoInfo.PacketLog) == 0 {
				continue
			}
			if hdoInfo.PacketLog[len(hdoInfo.PacketLog)-1].WithCoupon > volumeLimit {
				c.HdoInfo = append(c.HdoInfo, &coupon.HdoInfo{
					HdoServiceCode: hdoInfo.HdoServiceCode,
					CouponUse:      true,
				})
			}
		}

		for _, hduInfo := range packetLogInfo.HduInfo {
			if len(hduInfo.PacketLog) == 0 {
				continue
			}
			if hduInfo.PacketLog[len(hduInfo.PacketLog)-1].WithCoupon > volumeLimit {
				c.HduInfo = append(c.HduInfo, &coupon.HduInfo{
					HduServiceCode: hduInfo.HduServiceCode,
					CouponUse:      true,
				})
			}
		}

		for _, hdxInfo := range packetLogInfo.HdxInfo {
			if len(hdxInfo.PacketLog) == 0 {
				continue
			}
			if hdxInfo.PacketLog[len(hdxInfo.PacketLog)-1].WithCoupon > volumeLimit {
				c.HdxInfo = append(c.HdxInfo, &coupon.HdxInfo{
					HdxServiceCode: hdxInfo.HdxServiceCode,
					CouponUse:      true,
				})
			}
		}

		if len(c.HdoInfo) != 0 || len(c.HduInfo) != 0 || len(c.HdxInfo) != 0 {
			couponInfo = append(couponInfo, &c)
		}
	}

	if len(couponInfo) != 0 {
		body, statusCode, err := client.PutCoupon(accessToken, couponInfo)
		if err != nil {
			log.Println(err)
			return
		}
		if statusCode != http.StatusOK {
			log.Println(statusCode, body.ReturnCode)
			return
		}
		log.Println("switched off:")
		for _, c := range couponInfo {
			for _, hdoInfo := range c.HdoInfo {
				log.Printf("\t%s", hdoInfo.HdoServiceCode)
			}
			for _, hduInfo := range c.HduInfo {
				log.Printf("\t%s", hduInfo.HduServiceCode)
			}
			for _, hdxInfo := range c.HdxInfo {
				log.Printf("\t%s", hdxInfo.HdxServiceCode)
			}
		}
	}
}
