# fastnetmon_flow_spec_fragmentation
In this repo you can find callback which allows you to block fragmentaed traffic when FastNetMon dertects attack with zero ports

To start please configure FastNetMon [API](https://fastnetmon.com/docs-fnm-advanced/advanced-api/)

Then create file /etc/fastnetmon/fastnetmon_flow_spec_fragmentation.conf and put API credentials here:
```
{
  "api_user": "admin",
  "api_password": "put your password here",
  "api_host":  "127.0.0.1",
  "api_port": 10007
}
```

Then download binary file of integration and put it to /opt/fastnetmon_flow_spec_fragmentation and set chmod flag for it:
```
chmod +x /opt/fastnetmon_flow_spec_fragmentation
```


After that specify it on FastNetMon side as callback script:
```
sudo fcli set main notify_script_enabled enable
sudo fcli set main notify_script_format json
sudo fcli set main notify_script_path /opt/fastnetmon_flow_spec_fragmentation
sudo fcli commit
```

Try manually blocking following Flow Spec rule:
```
sudo fcli set flowspec  '{ "source_prefix": "4.0.0.0/32", "destination_prefix": "127.0.0.0/32", "destination_ports": [ 0 ], "source_ports": [ 0  ], "packet_lengths": [ 1500 ], "protocols": [ "udp" ], "action_type": "rate-limit", "action": { "rate": 1024 } }'
```

Then check that FastNetMon added supplementary Flow Spec announce:
```
sudo fcli show flowspec
{"action":{"rate":1024},"action_type":"rate-limit","destination_ports":[0],"destination_prefix":"127.0.0.0/32","packet_lengths":[1500],"protocols":["udp"],"source_ports":[0],"source_prefix":"4.0.0.0/32"} 30314e1f-d122-4f3f-8fcf-8cfbf3f7a427

{"action":{"rate":1024},"action_type":"rate-limit","destination_prefix":"127.0.0.0/32","fragmentation_flags":["is-fragment"],"protocols":["udp"],"source_prefix":"4.0.0.0/32"} c07ec922-76a6-40f0-accb-f7fcca2527c4
``` 

Then remove main announce;
```
sudo fcli delete flowspec 30314e1f-d122-4f3f-8fcf-8cfbf3f7a427
```

And check that FastNetMOn removed supplementary on too.

```
sudo fcli show flowspec
```
