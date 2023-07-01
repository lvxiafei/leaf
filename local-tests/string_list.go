package main

import (
	"log"
	"strings"
)

func main() {

	//eventOutput := "[03:27:13] [0xffff9f00e05ff000] [4026531992] 21904    2      eth0               fa:16:3e:0e:f4:e1  84     I_reply:1.1.1.1->172.17.0.2                            pkt_type=HOST iptables=[pf=PF_INET, table=filter hook=FORWARD verdict=ACCEPT cost=14.248µs]"
	eventOutput := "[08:03:39] [0xffff93b5af00a400] [4026532008] 38771 22 nil 02:00:00:00:00:00 00:49:4e:41:ed:03 0.0.0.0 0 15872 I_request:192.168.64.11->1.1.1.1 pkt_type=HOST func=ip_send_skb"
	item := strings.Fields(eventOutput)
	log.Printf("item slice: %+v", item)
	log.Printf("item slice len: %+v", len(item))
	log.Printf("item slice part: %+v", item[11:])
	log.Printf("item slice part: %+v", item[13:])
}
