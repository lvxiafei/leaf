package main

import (
	"leaf/tools"
	"leaf/util"
	"log"
	"strings"
)

func main() {
	var outputRequestItems []string
	var outputReplyItems []string

	eventOutput := "[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    nil            00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        15872    I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_send_skb\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    nil            00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        15872    I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_local_out\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    nil            00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    nil            00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_output\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    eth0           00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    eth0           00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_finish_output\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    eth0           00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_finish_output2\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    eth0           00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=dev_queue_xmit\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    eth0           00:00:00:03:00:00  00:00:00:00:00:00  192.168.6.16   0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=__dev_queue_xmit\n[08:31:21] [0xffff93b12527aa00] [4026534007] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=dev_hard_start_xmit\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=OTHERHOST func=netif_rx\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=OTHERHOST func=__netif_receive_skb\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=OTHERHOST func=packet_rcv\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=br_nf_pre_routing\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=iptable_mangle_hook\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      iptables=[pf=PF_INET table=mangle hook=PREROUTING verdict=ACCEPT cost=14.407µs]\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=iptable_nat_do_chain\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      iptables=[pf=PF_INET table=nat hook=PREROUTING verdict=ACCEPT cost=5.485µs]\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=br_nf_pre_routing_finish\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=br_handle_frame_finish\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    vethd39f4acf   ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=br_pass_frame_up\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=br_netif_receive_skb\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=__netif_receive_skb\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_rcv\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_rcv_finish\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=iptable_mangle_hook\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      iptables=[pf=PF_INET table=mangle hook=FORWARD verdict=ACCEPT cost=5.084µs]\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=iptable_filter_hook\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      iptables=[pf=PF_INET table=filter hook=FORWARD verdict=ACCEPT cost=12.318µs]\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    cni-podman0    ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=ip_output\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=iptable_mangle_hook\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      iptables=[pf=PF_INET table=mangle hook=POSTROUTING verdict=ACCEPT cost=4.608µs]\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      func=iptable_nat_do_chain\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:10.88.0.23->1.1.1.1                  pkt_type=HOST      iptables=[pf=PF_INET table=nat hook=POSTROUTING verdict=ACCEPT cost=8.985µs]\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=HOST      func=ip_finish_output\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           ba:f5:32:c3:f5:35  f6:81:f4:04:10:47  10.88.0.1      0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=HOST      func=ip_finish_output2\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           e8:61:1f:18:54:8d  00:00:5e:00:01:01  192.168.6.1    0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=HOST      func=dev_queue_xmit\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           e8:61:1f:18:54:8d  00:00:5e:00:01:01  192.168.6.1    0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=HOST      func=__dev_queue_xmit\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24332    ping         9    eth0           e8:61:1f:18:54:8d  00:00:5e:00:01:01  192.168.6.1    0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=HOST      func=dev_hard_start_xmit\n[08:31:21] [0xffff93b12527a000] [4026532008] 24332    ping         9    eth0           e8:61:1f:18:54:8d  00:00:5e:00:01:01  192.168.6.1    0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=OUTGOING  func=packet_rcv\n[08:31:21] [0xffff93b12527a000] [4026532008] 24332    ping         9    eth0           e8:61:1f:18:54:8d  00:00:5e:00:01:01  192.168.6.1    0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=OUTGOING  func=kfree_skbmem\n[08:31:21] [0xffff93b12527aa00] [4026532008] 24126    compile      17   eth0           e8:61:1f:18:54:8d  00:00:5e:00:01:01  192.168.6.1    0        84       I_request:192.168.6.11->1.1.1.1                pkt_type=HOST      func=skb_free_head\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      func=napi_gro_receive\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      func=__netif_receive_skb\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      func=packet_rcv\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      func=ip_rcv\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      func=iptable_mangle_hook\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->192.168.6.11                  pkt_type=HOST      iptables=[pf=PF_INET table=mangle hook=PREROUTING verdict=ACCEPT cost=11.281µs]\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_rcv_finish\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=iptable_mangle_hook\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      iptables=[pf=PF_INET table=mangle hook=FORWARD verdict=ACCEPT cost=4.872µs]\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=iptable_filter_hook\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      iptables=[pf=PF_INET table=filter hook=FORWARD verdict=ACCEPT cost=9.959µs]\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   eth0           6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_output\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=iptable_mangle_hook\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      iptables=[pf=PF_INET table=mangle hook=POSTROUTING verdict=ACCEPT cost=5.172µs]\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_finish_output\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    6c:e5:f7:71:4d:8c  e8:61:1f:18:54:8d  192.168.6.11   0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_finish_output2\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=dev_queue_xmit\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=__dev_queue_xmit\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=dev_hard_start_xmit\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=br_forward\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   cni-podman0    f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=__br_forward\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=br_forward_finish\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=br_nf_post_routing\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=dev_queue_xmit\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=__dev_queue_xmit\n[08:31:21] [0xffff93b749e05900] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=dev_hard_start_xmit\n[08:31:21] [0xffff93b749e05700] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=OUTGOING  func=packet_rcv\n[08:31:21] [0xffff93b749e05700] [4026532008] 23991    compile      18   vethd39f4acf   f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=OUTGOING  func=kfree_skbmem\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=netif_rx\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=__netif_receive_skb\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_rcv\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_rcv_finish\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ip_local_deliver\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=nf_hook_slow\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=icmp_rcv\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=ping_rcv\n[08:31:21] [0xffff93b749e05900] [4026534007] 23991    compile      18   eth0           f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=kfree_skbmem\n[08:31:21] [0xffff93b749e05700] [4026534007] 24332    ping         19   nil            f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=skb_free_head\n[08:31:21] [0xffff93b749e05700] [4026534007] 24332    ping         19   nil            f6:81:f4:04:10:47  ba:f5:32:c3:f5:35  10.88.0.23     0        84       I_reply:1.1.1.1->10.88.0.23                    pkt_type=HOST      func=kfree_skbmem"
	eventOutputList := strings.Split(eventOutput, "\n")

	for _, eventItem := range eventOutputList {
		item := strings.Fields(eventItem)
		if strings.Contains(item[12], "reply") {
			if len(item[14]) > 5 && item[14][:4] == "func" {
				outputReplyItems = append(outputReplyItems, item[14][5:])
			} else {
				outputReplyItems = append(outputReplyItems, item[14][:8])
			}
		} else {
			if len(item[14]) > 5 && item[14][:4] == "func" {
				outputRequestItems = append(outputRequestItems, item[14][5:])
			} else {
				outputRequestItems = append(outputRequestItems, item[14][:8])
			}
		}
	}

	j := tools.NewJsonData()
	j.AddSingleText(eventOutput)
	// Default Generate
	err := j.GenerateJsonFile(outputRequestItems, outputReplyItems, util.OutputJson)
	if err != nil {
		log.Printf("generate json data error: %v", err)
		return
	}
}