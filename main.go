package main

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"leaf/config"
	"leaf/model"
	"leaf/tools"
	"leaf/util"
	"leaf/versions"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unsafe"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -no-global-types -cc clang bpf ./bpf/skbtracer.c -- -D__TARGET_ARCH_x86 -I./bpf/headers -Wall
//go:embed scene-hook
var sceneHook embed.FS
var usage = `examples:
leaf                                      # trace all packets
leaf --proto=icmp -H 1.2.3.4 --icmpid 22  # trace icmp packet with addr=1.2.3.4 and icmpid=22
leaf --proto=tcp  -H 1.2.3.4 -P 22        # trace tcp  packet with addr=1.2.3.4:22
leaf --proto=udp  -H 1.2.3.4 -P 22        # trace udp  packet wich addr=1.2.3.4:22
leaf -t -T -p 1 -P 80 -H 127.0.0.1 --proto=tcp --icmpid=100 -N 10000
`
var RootCmd = cobra.Command{
	Use:   "leaf",
	Short: "network design and self-check through leaf",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Cfg.Parse(); err != nil {
			fmt.Println(err)
			return
		}
		config.RunGops()
		runEbpf()
	},
}

func main() {
	config.SetFlags(&RootCmd)
	cobra.CheckErr(RootCmd.Execute())
}

func runEbpf() {
	if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &unix.Rlimit{
		Cur: 4096,
		Max: 4096,
	}); err != nil {
		log.Fatalf("failed to set temporary rlimit: %s", err)
	}
	if err := unix.Setrlimit(unix.RLIMIT_MEMLOCK, &unix.Rlimit{
		Cur: unix.RLIM_INFINITY,
		Max: unix.RLIM_INFINITY,
	}); err != nil {
		log.Fatalf("Failed to set temporary rlimit: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	bpfSpec, err := loadBpf()
	if err != nil {
		log.Printf("Failed to load bpf spec: %v", err)
		return
	}

	if err := bpfSpec.RewriteConstants(map[string]interface{}{
		"CFG": config.GetBpfConfig(),
	}); err != nil {
		log.Printf("Failed to rewrite const for config: %v", err)
		return
	}

	var bpfObj bpfObjects
	if err := bpfSpec.LoadAndAssign(&bpfObj, &ebpf.CollectionOptions{
		Programs: ebpf.ProgramOptions{
			LogSize: ebpf.DefaultVerifierLogSize * 10,
		},
	}); err != nil {
		var ve *ebpf.VerifierError
		if errors.As(err, &ve) {
			log.Printf("Failed to load bpf obj: %v\n%-50v", err, ve)
		} else {
			log.Printf("Failed to load bpf obj: %v", err)
		}
		return
	}

	isHighVersion, err := versions.IsKernelVersionGte_5_16_0()
	if err != nil {
		log.Printf("Failed to check kernel version: %v", err)
		return
	}

	kIptDoTable := bpfObj.K_iptDoTable
	if !isHighVersion {
		kIptDoTable = bpfObj.OldK_iptDoTable
	}

	if kp, err := link.Kprobe("ipt_do_table", kIptDoTable, nil); err != nil {
		log.Printf("Failed to kprobe(ipt_do_table): %v", err)
		return
	} else {
		defer kp.Close()
		log.Printf("Attached kprobe(ipt_do_table)")
	}

	if krp, err := link.Kretprobe("ipt_do_table", bpfObj.KrIptDoTable, nil); err != nil {
		log.Printf("Failed to kretprobe(ipt_do_table): %v", err)
		return
	} else {
		defer krp.Close()
		log.Printf("Attached kretprobe(ipt_do_table)")
	}

	//if kp, err := link.Kprobe("ip6t_do_table", kIp6tDoTable, nil); err != nil {
	//	log.Printf("Failed to kprobe(ip6t_do_table): %v", err)
	//	return
	//} else {
	//	defer kp.Close()
	//	log.Printf("Attached kprobe(ip6t_do_table)")
	//}
	//
	//if krp, err := link.Kretprobe("ip6t_do_table", bpfObj.KrIp6tDoTable, nil); err != nil {
	//	log.Printf("Failed to kretprobe(ip6t_do_table): %v", err)
	//	return
	//} else {
	//	defer krp.Close()
	//	log.Printf("Attached kretprobe(ip6t_do_table)")
	//}

	//KprobeWrapper("ipt_do_table", bpfObj.OldFuncK_iptDoTable, nil)

	KprobeWrapper("nf_log_trace", bpfObj.K_nfLogTrace, nil)
	KprobeWrapper("kfree_skb_reason", bpfObj.K_kfreeSkbReason, nil)
	KprobeWrapper("ip_local_out", bpfObj.K_ipLocalOut, nil)
	KprobeWrapper("nf_hook_slow", bpfObj.K_nfHookSlow, nil)
	KprobeWrapper("kfree_skb", bpfObj.K_kfreeSkb, nil)
	KprobeWrapper("__kfree_skb", bpfObj.K___kfreeSkb, nil)
	KprobeWrapper("skb_free_head", bpfObj.K_skbFreeHead, nil)
	KprobeWrapper("kfree_skbmem", bpfObj.K_kfreeSkbmem, nil)
	KprobeWrapper("iptable_nat_do_chain", bpfObj.K_iptableNatDoChain, nil)
	KprobeWrapper("iptable_filter_hook", bpfObj.K_iptableFilterHook, nil)
	KprobeWrapper("iptable_mangle_hook", bpfObj.K_iptableMangleHook, nil)

	// netif rcv hook
	KprobeWrapper("netif_rx", bpfObj.K_netifRx, nil)
	KprobeWrapper("__netif_receive_skb", bpfObj.K___netifReceiveSkb, nil)
	KprobeWrapper("tpacket_rcv", bpfObj.K_tpacketRcv, nil)
	KprobeWrapper("packet_rcv", bpfObj.K_packetRcv, nil)
	KprobeWrapper("napi_gro_receive", bpfObj.K_napiGroReceive, nil)

	// netif send hook
	KprobeWrapper("dev_queue_xmit", bpfObj.K_devQueueXmit, nil)
	KprobeWrapper("__dev_queue_xmit", bpfObj.K___devQueueXmit, nil)
	KprobeWrapper("dev_hard_start_xmit", bpfObj.K_devHardStartXmit, nil)
	KprobeWrapper("loopback_xmit", bpfObj.K_loopbackXmit, nil)

	// br process hook
	KprobeWrapper("br_handle_frame_finish", bpfObj.K_brHandleFrameFinish, nil)
	KprobeWrapper("br_nf_pre_routing", bpfObj.K_brNfPreRouting, nil)
	KprobeWrapper("br_nf_pre_routing_finish", bpfObj.K_brNfPreRoutingFinish, nil)
	KprobeWrapper("br_pass_frame_up", bpfObj.K_brPassFrameUp, nil)
	KprobeWrapper("br_netif_receive_skb", bpfObj.K_brNetifReceiveSkb, nil)
	KprobeWrapper("br_forward", bpfObj.K_brForward, nil)
	KprobeWrapper("__br_forward", bpfObj.K___brForward, nil)
	KprobeWrapper("br_forward_finish", bpfObj.K_brForwardFinish, nil)
	KprobeWrapper("br_nf_forward_ip", bpfObj.K_brNfForwardIp, nil)
	KprobeWrapper("br_nf_forward_finish", bpfObj.K_brNfForwardFinish, nil)
	KprobeWrapper("br_nf_post_routing", bpfObj.K_brNfPostRouting, nil)
	KprobeWrapper("br_nf_dev_queue_xmit", bpfObj.K_brNfDevQueueXmit, nil)

	// vlan
	KprobeWrapper("skb_vlan_untag", bpfObj.K_skbVlanUntag, nil)
	KprobeWrapper("skb_vlan_push", bpfObj.K_skbVlanPush, nil)
	KprobeWrapper("skb_vlan_pop", bpfObj.K_skbVlanPop, nil)
	KprobeWrapper("vlan_dev_hard_start_xmit", bpfObj.K_vlanDevHardStartXmit, nil)
	KprobeWrapper("netif_receive_skb_core", bpfObj.K_netifReceiveSkbCore, nil)
	KprobeWrapper("vlan_do_receive", bpfObj.K_vlanDoReceive, nil)

	// l3 layer
	KprobeWrapper("ip_rcv", bpfObj.K_ipRcv, nil)
	KprobeWrapper("ip_rcv_finish", bpfObj.K_ipRcvFinish, nil)
	KprobeWrapper("ip_output", bpfObj.K_ipOutput, nil)
	KprobeWrapper("ip_finish_output", bpfObj.K_ipFinishOutput, nil)
	KprobeWrapper("ip_finish_output2", bpfObj.K_ipFinishOutput2, nil)
	KprobeWrapper("ip_send_skb", bpfObj.K_ipSendSkb, nil)
	KprobeWrapper("ip_queue_xmit", bpfObj.K_ipQueueXmit, nil)
	KprobeWrapper("icmp_rcv", bpfObj.K_icmpRcv, nil)
	KprobeWrapper("ip_local_deliver", bpfObj.K_ipLocalDeliver, nil)

	// l4 layer
	KprobeWrapper("ping_v4_sendmsg", bpfObj.K_ipSendSkb, nil)
	KprobeWrapper("ping_rcv", bpfObj.K_pingRcv, nil)

	// macvlan
	KprobeWrapper("macvlan_start_xmit", bpfObj.K_macvlanStartXmit, nil)
	//KretprobeWrapper("macvlan_start_xmit", bpfObj.KrMacvlanStartXmit, nil)

	KprobeWrapper("__neigh_event_send", bpfObj.K___neighEventSend, nil)
	KprobeWrapper("arp_solicit", bpfObj.K_arpSolicit, nil)
	KprobeWrapper("macvlan_broadcast_enqueue", bpfObj.K_macvlanBroadcastEnqueue, nil)

	rd, err := ringbuf.NewReader(bpfObj.EventsRingbuf)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %s", err)
	}

	var outputRequestItems []string
	var outputReplyItems []string
	var eventOutput string
	go func() {
		<-ctx.Done()
		_ = rd.Close()
		log.Println("Received signal, exiting program...")
	}()

	//log.Printf("Perf event size: %v", unsafe.Sizeof(perfEvent{})) // 208

	fmt.Printf("%-10s %-20s %-12s %-8s %-12s %-4s %-14s "+
		//"%-12s "+
		"%-18s %-18s %-14s %-8s %-8s "+
		"%-46s %s\n",
		"TIME", "SKB", "NETWORK_NS", "PID", "COMM", "CPU", "INTERFACE",
		//"EthProto",
		"SRC_MAC", "DEST_MAC", "NEXTHOP", "VLAN_ID", "IP_LEN",
		"PKT_INFO", "TRACE_INFO")

	var event model.PerfEvent
	var ipt model.IptablesInfo
	var trace model.IptablesTrace
	//indexIptablesInfo := unsafe.Sizeof(perfEvent{})
	indexIptablesInfo := unsafe.Sizeof(model.PerfEvent{}) - unsafe.Sizeof(model.IptablesInfo{}) - 4

	forever := config.Cfg.CatchCount == 0
	for n := config.Cfg.CatchCount; forever || n > 0; n-- {

		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("received signal, exiting..")
				if config.Cfg.Draw {
					j := tools.NewJsonData()
					j.AddSingleText(eventOutput)

					// Node2Eip Define
					if config.Cfg.Scene == util.Node2Eip {
						j.OpenSceneHookFsToAppend(sceneHook, util.Node2EipDefine)
					}
					// Node2Eip Diff
					if config.Cfg.Scene == util.Node2Eip && config.Cfg.Mode == util.Diff {
						if event.L4Proto == model.IpprotoICMP {
							j.OpenSceneHookFsToAppend(sceneHook, util.Node2EipDiffICMP)
						} else if event.L4Proto == model.IpprotoTCP {
							j.OpenSceneHookFsToAppend(sceneHook, util.Node2EipDiffTCP)
						} else if event.L4Proto == model.IpprotoUDP {
							j.OpenSceneHookFsToAppend(sceneHook, util.Node2EipDiffUDP)
						} else {
							j.OpenSceneHookFsToAppend(sceneHook, util.Node2EipDiff)
						}
					}

					// Pod2PodDifferentNode Define
					if config.Cfg.Scene == util.Pod2PodDifferentNode {
						j.OpenSceneHookFsToAppend(sceneHook, util.Pod2PodDifferentNodeDefine)
					}

					// Default Generate
					err := j.GenerateJsonFile(outputRequestItems, outputReplyItems, util.OutputJson)
					if err != nil {
						log.Printf("generate json data error: %v", err)
						return
					}

					// copy
					//err = tools.Copy2Clipboard()
					//if err != nil {
					//	log.Printf("Copy2Clipboard error: %v", err)
					//	return
					//}

					// upload
					err, outStr := tools.UploadFile()
					if err != nil {
						log.Printf("upload data error: %v", err)
						return
					}
					log.Printf("upload data:\n%v", outStr)
				}
				return
			}
			log.Printf("reading from reader: %s", err)
			continue
		}
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("Failed to parse perf event: %v", err)
			continue
		}
		if config.Cfg.Debug.Default || config.Cfg.Debug.IndexIptablesInfo {
			log.Printf("record.RawSample: %+v", record.RawSample)
			log.Printf("event: %+v", event)
			log.Printf("event.iptablesInfo: %+v", event.IptablesInfo)
		}

		if event.Flags&model.RouteEventIptable == model.RouteEventIptable {
			if err := binary.Read(bytes.NewReader(record.RawSample[indexIptablesInfo:]), binary.LittleEndian, &ipt); err != nil {
				log.Printf("Failed to parse iptables info: %v", err)
				continue
			}
			if config.Cfg.Debug.Default {
				log.Printf("iptablesInfo: %+v", ipt)
			}
		} else if event.Flags&model.RouteEventIptablesTrace == model.RouteEventIptablesTrace {
			if err := binary.Read(bytes.NewReader(record.RawSample[indexIptablesInfo:]), binary.LittleEndian, &trace); err != nil {
				log.Printf("Failed to parse iptables trace: %v", err)
				continue
			}
			if config.Cfg.Debug.Default {
				log.Printf("iptablesTrace: %+v", trace)
			}
		}

		eventItem := event.Output(&ipt, &trace)
		//eventItem := event.output(&event.iptablesInfo, &trace)

		fmt.Println(eventItem)

		eventOutput = eventOutput + eventItem + "\n"

		if config.Cfg.Draw {
			item := strings.Fields(eventItem)
			if config.Cfg.Debug.Default {
				log.Printf("item: %+v", item)
				log.Printf("item len: %+v", len(item))
			}

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

		select {
		case <-ctx.Done():
			log.Println("ctx done ...")
			return
		default:
		}
	}
}

func KprobeWrapper(symbol string, prog *ebpf.Program, opts *link.KprobeOptions) {

	if len(config.Cfg.File) > 0 && !config.Cfg.Hook[symbol] {
		return
	}

	if _, err := link.Kprobe(symbol, prog, opts); err != nil {
		log.Printf("Failed to kprobe(%s): %v", symbol, err)
		return
	} else {
		log.Printf("Attached kprobe(%s)", symbol)
	}
}

func KretprobeWrapper(symbol string, prog *ebpf.Program, opts *link.KprobeOptions) {

	if len(config.Cfg.File) > 0 && !config.Cfg.Hook[symbol] {
		return
	}

	if _, err := link.Kretprobe(symbol, prog, opts); err != nil {
		log.Printf("Failed to kretprobe(%s): %v", symbol, err)
		return
	} else {
		log.Printf("Attached kretprobe(%s)", symbol)
	}
}
