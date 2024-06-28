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
