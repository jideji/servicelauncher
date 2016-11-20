# Service Launcher

**Note: This project is currently in development and isn't even in ALPHA state.**

Run/restart/stop software locally for development.

It expects a configuration file in the home directory ($HOME/.slcfg).
Each service is configured with the format:

```
service.<servicename>.pattern = <regex for finding command in ps output>
service.<servicename>.command = <command for starting service>
```

e.g.

```
service.httpserver.pattern = Python -m SimpleHTTPServer 8080
service.httpserver.command = python -m SimpleHTTPServer 8080

service.ncserver.pattern = ^nc -l 8081$
service.ncserver.command = nc -l 8081
```
