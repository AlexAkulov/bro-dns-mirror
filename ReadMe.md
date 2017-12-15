# Mirror dns queries form bro-logs to Cisco Umbrella

```
ssh sd2-bro1.example.com "tail -f /opt/bro/logs/current/dns.log | cut -d$'\t' -f 10,14" | ./bro-dns-mirror
```
