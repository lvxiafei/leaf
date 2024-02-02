#include "skbtracer.h"
#include "bpf_kprobe_args.h"
/*
 * Common tracepoint handler. Detect IPv4/IPv6 and
 * emit event with address, interface and namespace.
 */
static __inline bool
do_trace_skb(struct event_t *event,
    struct pt_regs *ctx,
    struct sk_buff *skb)
{
    unsigned char *l3_header;
    u8 ip_version, l4_proto;

    event->flags |= SKBTRACER_EVENT_IF;
    set_event_info(skb, event);
    set_pkt_info(skb, &event->pkt_info);
    set_ether_info(skb, &event->l2_info);

    l3_header = get_l3_header(skb);
    ip_version = get_ip_version(l3_header);
    if (ip_version == 4) {
        event->l2_info.l3_proto = ETH_P_IP;
        set_ipv4_info(skb, &event->l3_info);
    } else if (ip_version == 6) {
        event->l2_info.l3_proto = ETH_P_IPV6;
        set_ipv6_info(skb, &event->l3_info);
    } else {
        return false;
    }

    l4_proto = event->l3_info.l4_proto;
    if (l4_proto == IPPROTO_TCP) {
        set_tcp_info(skb, &event->l4_info);
    } else if (l4_proto == IPPROTO_UDP) {
        set_udp_info(skb, &event->l4_info);
    } else if (l4_proto == IPPROTO_ICMP || l4_proto == IPPROTO_ICMPV6) {
        set_icmp_info(skb, &event->icmp_info);
    } else {
        return false;
    }

    return true;
}

static __inline int
do_trace(struct pt_regs *ctx,
    struct sk_buff *skb,
    const char *func_name)
{

    u64 pid_tgid;
    pid_tgid = bpf_get_current_pid_tgid();

    if (filter_pid(pid_tgid>>32) || filter_netns(skb) || filter_l3_and_l4_info(skb))
        return false;

    struct event_t *event = GET_EVENT_BUF();
    if (!event)
        return BPF_OK;

    if (!do_trace_skb(event, ctx, skb)) return 0;

//    if (!filter_callstack(cfg))
//        set_callstack(event, ctx);

//    bpf_strncpy(event->func_name, func_name, FUNCNAME_MAX_LEN);
    bpf_probe_read(&event->func_name, sizeof(event->func_name), func_name);

    bpf_ringbuf_output(&events_ringbuf, event, sizeof(*event), 0);

    return 0;
}

static __noinline int
__ipt_do_table_in(struct pt_regs *ctx,
    struct sk_buff *skb,
    const struct nf_hook_state *state,
    struct xt_table *table)
{
    u64 pid_tgid;
    pid_tgid = bpf_get_current_pid_tgid();

    if (filter_pid(pid_tgid>>32) || filter_netns(skb) || filter_l3_and_l4_info(skb))
        return false;

    struct ipt_do_table_args args = {
        .skb = skb,
        .state = state,
        .table = table,
    };

    args.start_ns = bpf_ktime_get_ns();
    bpf_map_update_elem(&skbtracer_ipt, &pid_tgid, &args, BPF_ANY);

    return BPF_OK;
};

static __noinline int
__ipt_do_table_out(struct pt_regs *ctx, uint verdict)
{
    u64 pid_tgid;
    u64 ipt_delay;
    struct ipt_do_table_args *args;

    pid_tgid = bpf_get_current_pid_tgid();
    args = bpf_map_lookup_elem(&skbtracer_ipt, &pid_tgid);
    if (args == NULL)
        return BPF_OK;

    bpf_map_delete_elem(&skbtracer_ipt, &pid_tgid);

    struct event_t *event = GET_EVENT_BUF();
    if (!event)
        return BPF_OK;

    if (!do_trace_skb(event, ctx, args->skb))
        return BPF_OK;

    event->flags |= SKBTRACER_EVENT_IPTABLE;

    ipt_delay = bpf_ktime_get_ns() - args->start_ns;
    set_iptables_info(args->table, args->state, (u32)verdict, ipt_delay, &event->ipt_info);

    bpf_ringbuf_output(&events_ringbuf, event, sizeof(*event), 0);

    return BPF_OK;
}

static __inline int
__ipt_do_table_trace(struct pt_regs *ctx,
    u8 pf,
    unsigned int hooknum,
    struct sk_buff *skb,
    struct net_device *in,
    struct net_device *out,
    char *tablename,
    char *chainname,
    unsigned int rulenum)
{
    u64 pid_tgid = bpf_get_current_pid_tgid();
    void *val = bpf_map_lookup_elem(&skbtracer_ipt, &pid_tgid);
    if (!val)
        return BPF_OK;

    struct event_t *event = GET_EVENT_BUF();
    if (!event)
        return BPF_OK;

    __builtin_memset(event, 0, sizeof(*event));

    if (!do_trace_skb(event, ctx, skb))
        return BPF_OK;

    event->flags |= SKBTRACER_EVENT_IPTABLES_TRACE;

    struct iptables_trace_t *trace = &event->trace_info;

    read_dev_name((char *)&trace->in, in);
    read_dev_name((char *)&trace->out, out);
    bpf_probe_read_kernel_str(&trace->tablename, XT_TABLE_MAXNAMELEN, tablename);
    bpf_probe_read_kernel_str(&trace->chainname, XT_TABLE_MAXNAMELEN, chainname);
    trace->rulenum = (u32)rulenum;
    trace->hooknum = (u32)hooknum;
    trace->pf = pf;

    bpf_ringbuf_output(&events_ringbuf, event, sizeof(*event), 0);

    return BPF_OK;
}

// >= 5.16

SEC("kprobe/ipt_do_table")
int BPF_KPROBE(k_ipt_do_table, struct xt_table *table,
    struct sk_buff *skb,
    const struct nf_hook_state *state)
{
    return __ipt_do_table_in(ctx, skb, state, table);
};

// < 5.16

SEC("kprobe/ipt_do_table")
int BPF_KPROBE(old_k_ipt_do_table, struct sk_buff *skb,
    const struct nf_hook_state *state,
    struct xt_table *table)
{
    return __ipt_do_table_in(ctx, skb, state, table);
}

// only func
SEC("kprobe/ipt_do_table")
int BPF_KPROBE(old_func_k_ipt_do_table)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "ipt_do_table");
}

SEC("kretprobe/ipt_do_table")
int BPF_KRETPROBE(kr_ipt_do_table, uint ret)
{
    return __ipt_do_table_out(ctx, ret);
}

// SEC("kprobe/ip6t_do_table")
// int BPF_KPROBE(k_ip6t_do_table, void *priv, struct sk_buff *skb,
//     const struct nf_hook_state *state)
// {
//     struct xt_table *table = (struct xt_table *)priv;
//     return __ipt_do_table_in(ctx, skb, state, table);
// };

// SEC("kretprobe/ip6t_do_table")
// int BPF_KRETPROBE(kr_ip6t_do_table, uint ret)
// {
//     return __ipt_do_table_out(ctx, ret);
// }


SEC("kprobe/nf_log_trace")
int BPF_KPROBE(k_nf_log_trace, struct net *net, u_int8_t pf, unsigned int hooknum,
    struct sk_buff *skb, struct net_device *in)
{
    struct net_device *out;
    char *tablename;
    char *chainname;
    unsigned int rulenum;

    out = (typeof(out))(void *)regs_get_nth_argument(ctx, 5);
    tablename = (typeof(tablename))(void *)regs_get_nth_argument(ctx, 8);
    chainname = (typeof(chainname))(void *)regs_get_nth_argument(ctx, 9);
    rulenum = (typeof(rulenum))regs_get_nth_argument(ctx, 11);

    return __ipt_do_table_trace(ctx, pf, hooknum, skb, in, out, tablename,
        chainname, rulenum);
}



/*
 * netif rcv hook:
 * 1) int netif_rx(struct sk_buff *skb)
 * 2) int __netif_receive_skb(struct sk_buff *skb)
 * 3) gro_result_t napi_gro_receive(struct napi_struct *napi, struct sk_buff *skb)
 * 4) ...
 */

SEC("kprobe/netif_rx")
int BPF_KPROBE(k_netif_rx) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "netif_rx");
}

SEC("kprobe/__netif_receive_skb")
int BPF_KPROBE(k___netif_receive_skb) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "__netif_receive_skb");
}

SEC("kprobe/tpacket_rcv")
int BPF_KPROBE(k_tpacket_rcv) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "tpacket_rcv");
}

// int packet_rcv(struct sk_buff *skb, struct net_device *dev, struct packet_type *pt, struct net_device *orig_dev);
SEC("kprobe/packet_rcv")
int BPF_KPROBE(k_packet_rcv) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "packet_rcv");
}

SEC("kprobe/napi_gro_receive")
int BPF_KPROBE(k_napi_gro_receive) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "napi_gro_receive");
}

/*
 * netif send hook:
 * 1) int __dev_queue_xmit(struct sk_buff *skb, struct net_device *sb_dev)
 * 2) struct sk_buff *dev_hard_start_xmit(struct sk_buff *skb, struct net_device *dev, struct netdev_queue *txq, int *ret);
 * 3) netdev_tx_t loopback_xmit(struct sk_buff *skb, struct net_device *dev);
 * 4) ...
 */

// int dev_queue_xmit(struct sk_buff *skb);
SEC("kprobe/dev_queue_xmit")
int BPF_KPROBE(k_dev_queue_xmit) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "dev_queue_xmit");
}

SEC("kprobe/__dev_queue_xmit")
int BPF_KPROBE(k___dev_queue_xmit) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "__dev_queue_xmit");
}

// capture dev send packet
SEC("kprobe/dev_hard_start_xmit")
int BPF_KPROBE(k_dev_hard_start_xmit)
{
   struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
   return do_trace(ctx, skb, "dev_hard_start_xmit");
}

SEC("kprobe/loopback_xmit")
int BPF_KPROBE(k_loopback_xmit)
{
   struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
   return do_trace(ctx, skb, "loopback_xmit");
}

/*
 * br process hook:
 * 1) rx_handler_result_t br_handle_frame(struct sk_buff **pskb)
 * 2) int br_handle_frame_finish(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 3) unsigned int br_nf_pre_routing(void *priv, struct sk_buff *skb, const struct nf_hook_state *state)
 * 4) int br_nf_pre_routing_finish(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 5) int br_pass_frame_up(struct sk_buff *skb)
 * 6) int br_netif_receive_skb(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 7) void br_forward(const struct net_bridge_port *to, struct sk_buff *skb, bool local_rcv, bool local_orig)
 * 8) int br_forward_finish(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 9) unsigned int br_nf_forward_ip(void *priv,struct sk_buff *skb,const struct nf_hook_state *state)
 * 10)int br_nf_forward_finish(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 11)unsigned int br_nf_post_routing(void *priv,struct sk_buff *skb,const struct nf_hook_state *state)
 * 12)int br_nf_dev_queue_xmit(struct net *net, struct sock *sk, struct sk_buff *skb)
 */

SEC("kprobe/br_handle_frame_finish")
int BPF_KPROBE(k_br_handle_frame_finish) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "br_handle_frame_finish");
}

SEC("kprobe/br_nf_pre_routing")
int BPF_KPROBE(k_br_nf_pre_routing) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "br_nf_pre_routing");
}

SEC("kprobe/br_nf_pre_routing_finish")
int BPF_KPROBE(k_br_nf_pre_routing_finish) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "br_nf_pre_routing_finish");
}

SEC("kprobe/br_pass_frame_up")
int BPF_KPROBE(k_br_pass_frame_up) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "br_pass_frame_up");
}

SEC("kprobe/br_netif_receive_skb")
int BPF_KPROBE(k_br_netif_receive_skb) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "br_netif_receive_skb");
}

SEC("kprobe/br_forward")
int BPF_KPROBE(k_br_forward) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "br_forward");
}

SEC("kprobe/__br_forward")
int BPF_KPROBE(k___br_forward) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "__br_forward");
}

SEC("kprobe/br_forward_finish")
int BPF_KPROBE(k_br_forward_finish) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "br_forward_finish");
}

SEC("kprobe/br_nf_forward_ip")
int BPF_KPROBE(k_br_nf_forward_ip) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "br_nf_forward_ip");
}

SEC("kprobe/br_nf_forward_finish")
int BPF_KPROBE(k_br_nf_forward_finish) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "br_nf_forward_finish");
}

SEC("kprobe/br_nf_post_routing")
int BPF_KPROBE(k_br_nf_post_routing) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "br_nf_post_routing");
}

SEC("kprobe/br_nf_dev_queue_xmit")
int BPF_KPROBE(k_br_nf_dev_queue_xmit) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "br_nf_dev_queue_xmit");
}

/*
 * ip layer:
 * 1) int ip_rcv(struct sk_buff *skb, struct net_device *dev, struct packet_type *pt, struct net_device *orig_dev)
 * 2) int ip_rcv_finish(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 3) int ip_output(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 4) int ip_finish_output(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 5) int ip_finish_output2(struct net *net, struct sock *sk, struct sk_buff *skb)
 * 6) ...
 */

SEC("kprobe/ip_rcv")
int BPF_KPROBE(k_ip_rcv) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "ip_rcv");
}

SEC("kprobe/ip_rcv_finish")
int BPF_KPROBE(k_ip_rcv_finish) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "ip_rcv_finish");
}

SEC("kprobe/ip_output")
int BPF_KPROBE(k_ip_output) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "ip_output");
}

SEC("kprobe/ip_finish_output")
int BPF_KPROBE(k_ip_finish_output) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "ip_finish_output");
}

SEC("kprobe/ip_finish_output2")
int BPF_KPROBE(k_ip_finish_output2) {
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "ip_finish_output2");
}

// ip_send_skb, first
// int ip_send_skb(struct net *net, struct sk_buff *skb);
SEC("kprobe/ip_send_skb")
int BPF_KPROBE(k_ip_send_skb)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "ip_send_skb");
}

// int ip_queue_xmit(struct sock *sk, struct sk_buff *skb, struct flowi *fl);
SEC("kprobe/ip_queue_xmit")
int BPF_KPROBE(k_ip_queue_xmit)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "ip_queue_xmit");
}

// int ping_v4_sendmsg(struct sock *sk, struct msghdr *msg, size_t len);
SEC("kprobe/ping_v4_sendmsg")
int BPF_KPROBE(k_ping_v4_sendmsg)
{
   return do_trace(ctx, NULL, "ping_v4_sendmsg");
}

//


// int icmp_rcv(struct sk_buff *skb);
SEC("kprobe/icmp_rcv")
int BPF_KPROBE(k_icmp_rcv)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "icmp_rcv");
}

// bool ping_rcv(struct sk_buff *skb);
SEC("kprobe/ping_rcv")
int BPF_KPROBE(k_ping_rcv)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "ping_rcv");
}

// vlan
//struct sk_buff *skb_vlan_untag(struct sk_buff *skb);
SEC("kprobe/skb_vlan_untag")
int BPF_KPROBE(k_skb_vlan_untag)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "skb_vlan_untag");
}

//int skb_vlan_push(struct sk_buff *skb, __be16 vlan_proto, u16 vlan_tci);
SEC("kprobe/skb_vlan_push")
int BPF_KPROBE(k_skb_vlan_push)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "skb_vlan_push");
}

//int skb_vlan_pop(struct sk_buff *skb);
SEC("kprobe/skb_vlan_pop")
int BPF_KPROBE(k_skb_vlan_pop)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "skb_vlan_pop");
}

//netdev_tx_t vlan_dev_hard_start_xmit(struct sk_buff *skb, struct net_device *dev);
SEC("kprobe/vlan_dev_hard_start_xmit")
int BPF_KPROBE(k_vlan_dev_hard_start_xmit)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "vlan_dev_hard_start_xmit");
}

//int netif_receive_skb_core(struct sk_buff *skb);
SEC("kprobe/netif_receive_skb_core")
int BPF_KPROBE(k_netif_receive_skb_core)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "netif_receive_skb_core");
}

//bool vlan_do_receive(struct sk_buff *skbp);
SEC("kprobe/vlan_do_receive")
int BPF_KPROBE(k_vlan_do_receive)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "vlan_do_receive");
}


// iptables
//void kfree_skb_reason(struct sk_buff *skb, enum skb_drop_reason reason);
SEC("kprobe/kfree_skb_reason")
int BPF_KPROBE(k_kfree_skb_reason)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "kfree_skb_reason");
}

//int ip_local_out(struct net *net, struct sock *sk, struct sk_buff *skb);
SEC("kprobe/ip_local_out")
int BPF_KPROBE(k_ip_local_out)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "ip_local_out");
}

//int nf_hook_slow(struct sk_buff *skb, struct nf_hook_state *state, struct nf_hook_entries *e, unsigned int s);
SEC("kprobe/nf_hook_slow")
int BPF_KPROBE(k_nf_hook_slow)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "nf_hook_slow");
}

//void kfree_skb(struct sk_buff *skb)
SEC("kprobe/kfree_skb")
int BPF_KPROBE(k_kfree_skb)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "kfree_skb");
}
//void __kfree_skb(struct sk_buff *skb);
SEC("kprobe/__kfree_skb")
int BPF_KPROBE(k___kfree_skb)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "__kfree_skb");
}

//void skb_free_head(struct sk_buff *skb);
SEC("kprobe/skb_free_head")
int BPF_KPROBE(k_skb_free_head)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "skb_free_head");
}

//void kfree_skbmem(struct sk_buff *skb);
SEC("kprobe/kfree_skbmem")
int BPF_KPROBE(k_kfree_skbmem)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "kfree_skbmem");
}
// 4.18
//static unsigned int iptable_nat_do_chain(void *priv, struct sk_buff *skb, const struct nf_hook_state *state)
SEC("kprobe/iptable_nat_do_chain")
int BPF_KPROBE(k_iptable_nat_do_chain)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "iptable_nat_do_chain");
}

//static unsigned int iptable_filter_hook(void *priv, struct sk_buff *skb, const struct nf_hook_state *state)
SEC("kprobe/iptable_filter_hook")
int BPF_KPROBE(k_iptable_filter_hook)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "iptable_filter_hook");
}

//static unsigned int iptable_mangle_hook(void *priv, struct sk_buff *skb, const struct nf_hook_state *state)
SEC("kprobe/iptable_mangle_hook")
int BPF_KPROBE(k_iptable_mangle_hook)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "iptable_mangle_hook");
}

// 报文是发送给本机的话那么就会直接到 ip_local_deliver() 进行处理
//int ip_local_deliver(struct sk_buff *skb)
SEC("kprobe/ip_local_deliver")
int BPF_KPROBE(k_ip_local_deliver)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "ip_local_deliver");
}

// macvlan
//netdev_tx_t macvlan_start_xmit(struct sk_buff *skb, struct net_device *dev);
SEC("kprobe/macvlan_start_xmit")
int BPF_KPROBE(k_macvlan_start_xmit)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "macvlan_start_xmit");
}

SEC("kretprobe/macvlan_start_xmit")
int BPF_KRETPROBE(kr_macvlan_start_xmit)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM1(ctx);
    return do_trace(ctx, skb, "kr_macvlan_start_xmit");
}

//int __neigh_event_send(struct neighbour *neigh, struct sk_buff *skb);
SEC("kprobe/__neigh_event_send")
int BPF_KPROBE(k___neigh_event_send)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "__neigh_event_send");
}

//void arp_solicit(struct neighbour *neigh, struct sk_buff *skb);
SEC("kprobe/arp_solicit")
int BPF_KPROBE(k_arp_solicit)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM2(ctx);
    return do_trace(ctx, skb, "arp_solicit");
}

//void macvlan_broadcast_enqueue(struct macvlan_port *port, struct macvlan_dev *src, struct sk_buff *skb);
SEC("kprobe/macvlan_broadcast_enqueue")
int BPF_KPROBE(k_macvlan_broadcast_enqueue)
{
    struct sk_buff *skb = (struct sk_buff *)PT_REGS_PARM3(ctx);
    return do_trace(ctx, skb, "macvlan_broadcast_enqueue");
}

// tcp
//int tcp_sendmsg(struct sock *sk, struct msghdr *msg, size_t size);


// ipsec
//int xfrm4_output(struct net *net, struct sock *sk, struct sk_buff *skb);
