package main

import (
	"leaf/tools"
	"strconv"
)

func main() {
	items := []string{"a", "b", "c", "d"}
	generateFile := "generate_base.json"
	j := tools.NewJsonData()
	//text := "TIME       SKB                  NETWORK_NS   PID      CPU    INTERFACE          DEST_MAC           IP_LEN PKT_INFO                                               IPTABLES_INFO\n[00:00:56] [0xffff95398a31a200] [4026531840] 2347     0      nil                61:6e:37:38:78:78  84     I_request:10.0.2.15->8.8.8.8                           pkt_type=HOST iptables=[pf=PF_INET, table=nat hook=OUTPUT verdict=ACCEPT]\n[00:00:56] [0xffff95398a31a200] [4026531840] 2347     0      nil                61:6e:37:38:78:78  84     I_request:10.0.2.15->8.8.8.8                           pkt_type=HOST iptables=[pf=PF_INET, table=filter hook=OUTPUT verdict=ACCEPT]\n[00:00:56] [0xffff95398a31a200] [4026531840] 2347     0      enp0s3             61:6e:37:38:78:78  84     I_request:10.0.2.15->8.8.8.8                           pkt_type=HOST iptables=[pf=PF_INET, table=nat hook=POSTROUTING verdict=ACCEPT]\n[00:00:56] [0xffff953990ac5500] [4026531840] 0        3      enp0s3             08:00:27:ff:1e:ab  84     I_reply:8.8.8.8->10.0.2.15                             pkt_type=HOST iptables=[pf=PF_INET, table=filter hook=INPUT verdict=ACCEPT]\n[00:00:57] [0xffff95398a31ac00] [4026531840] 2347     0      nil                00:00:00:00:00:00  84     I_request:10.0.2.15->8.8.8.8                           pkt_type=HOST iptables=[pf=PF_INET, table=filter hook=OUTPUT verdict=ACCEPT]\n[00:00:57] [0xffff953990ac5100] [4026531840] 0        3      enp0s3             08:00:27:ff:1e:ab  84     I_reply:8.8.8.8->10.0.2.15                             pkt_type=HOST iptables=[pf=PF_INET, table=filter hook=INPUT verdict=ACCEPT]"
	//j.AddSingleText(text)

	nextRectangle := j.AddHeadRectangle(items[0])
	for i, item := range items[1:] {
		nextRectangle = j.AddDownArrowRectangleWithText(j, nextRectangle, item, strconv.Itoa(i+1))
	}

	//nextRectangle = j.AddHeadRectangleRightX(items[0], 400)
	nextRectangle = j.AddHeadRectangleLastRightX(nextRectangle, items[0], 400)
	for i, item := range items[1:] {
		nextRectangle = j.AddUpArrowRectangleWithTextRightX(j, nextRectangle, item, strconv.Itoa(i+1), 400)
	}

	//hook
	j.OpenToAppend("scene_hook_base/node2eip_define.json")

	j.WriteToJson(generateFile)
}
