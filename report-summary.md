# Firewall Escape Test - Run 21858986051

**Date**: 2026-02-10
**Status**: SANDBOX SECURE ✅
**Techniques Tested**: 30
**Novelty Rate**: 93.3%
**Escapes Found**: 0

## Key Findings
- IPv6 completely disabled (network unreachable) - prevents all IPv6-based bypasses
- SCTP protocol sockets can be created but connections are NAT'd to Squid and blocked
- Multicast works locally but has no external reach (no bypass)
- Netlink sockets allow route introspection but not modification or bypass
- All capabilities requiring NET_RAW/NET_ADMIN/SYS_ADMIN confirmed dropped
- /proc filesystem properly restricted: /proc/sys/net/, /proc/1/root/, /proc/net/arp all blocked
- TTL manipulation ineffective - NAT operates at first hop (Squid at 172.30.0.10)
- Environment variable manipulation irrelevant - iptables NAT enforces proxy at kernel level
- ARP cache inspection reveals gateway (172.30.0.1) and Squid (172.30.0.10) but provides no bypass
- All non-HTTP protocols (Gopher, RTSP, SCTP, NFS, SMB) properly blocked

## Architecture Strengths (Confirmed via Source Code Analysis + Testing)
1. **Triple-layer defense**: iptables NAT → Squid proxy → Host-level iptables
2. **Universal traffic redirection**: ALL TCP/UDP redirected to Squid (172.30.0.10:3128)
3. **DNS hardening**: Only trusted resolvers (8.8.8.8, 8.8.4.4, 127.0.0.11)
4. **Localhost bypass**: Intentionally allows 127.0.0.0/8 for MCP servers
5. **Host gateway bypass**: host.docker.internal bypasses Squid when AWF_ENABLE_HOST_ACCESS=1 (for MCP gateway)
6. **Capabilities properly dropped**: NET_RAW, NET_ADMIN, SYS_PTRACE, SYS_ADMIN, BPF all disabled
7. **Seccomp mode 2**: Syscall filtering active
8. **IPv6 disabled**: Network unreachable prevents IPv6 attack vectors

## Reconnaissance Intel Gathered This Run
- Gateway web service: 172.30.0.1 (MAC: 3e:df:34:29:17:f1)
- Squid proxy: 172.30.0.10:3128 (version 6.13, MAC: 5a:4c:5e:5d:db:62)
- Docker socket: /var/run/docker.sock (protected - connection refused)
- Network namespace: Isolated
- SCTP protocol: Available but NAT'd to Squid
- Multicast: Local-only (224.0.0.1)
- Netlink: Route table readable (default via 172.30.0.1)
- bpftool: Present but access denied
- IPv6: Completely unavailable

## Novel Techniques Introduced (93.3% - 28/30 NEW)
1. ✅ IPv6 DNS queries to Cloudflare IPv6 DNS
2. ✅ IPv6 direct HTTP connections
3. ✅ QUIC/HTTP3 protocol testing
4. ✅ ICMP data payload tunneling
5. ✅ eBPF/bpftool exploitation attempts
6. ✅ Network namespace creation (unshare)
7. ✅ HTTP/0.9 legacy protocol
8. ✅ Gopher protocol (pre-HTTP)
9. ✅ RTSP streaming protocol
10. ✅ SCTP protocol (Stream Control Transmission)
11. ✅ Multicast group joining (IP_ADD_MEMBERSHIP)
12. ✅ Netlink socket access (AF_NETLINK)
13. ✅ AF_PACKET raw socket attempts
14. ✅ Network interface manipulation (MTU)
15. ✅ DNS amplification (large buffer queries)
16. ✅ TTL manipulation (IP_TTL=1)
17. ✅ IP fragmentation attacks
18. ✅ UDP hole punching
19. ✅ ARP cache inspection
20. ✅ Conntrack manipulation attempts
21. ✅ Systematic PAC environment clearing
22. ✅ IPv6 Teredo tunneling (miredo)
23. ✅ /proc/sys/net kernel parameters
24. ✅ Container escape via /proc/1/root
25. ✅ Socat TCP relay
26. ✅ NFS/SMB file sharing protocols
27. ✅ Node.js dgram UDP module
28. ✅ Perl IO::Socket::SSL

## Recommendations
- Continue monitoring host.docker.internal bypass behavior
- Regular Squid version updates (currently 6.13)
- If IPv6 enabled in future, ensure ip6tables rules comprehensive
- Continue diverse security testing with >80% novel techniques per run

## Cumulative Statistics
- Total techniques: 568 (23 runs)
- Historical escapes: 1 (patched in AWF v0.9.1)
- Success rate: 0.18% (1/568)
- Last 538 techniques: All blocked
