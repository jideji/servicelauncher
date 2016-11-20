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

### Ideas

 - Command line tool
 - Terminal UI
 - Web UI with the same controls
 - Ability to see logs
   - Command line:
     - cat
     - tail
   - UI:
     - Highlight failures (based on regex?)
       - Scroll to failure
     - Allow for automatically switching log file displayed when an exception happens
 - Support for labeling services (e.g. materializers)
   - Make sure all services with a certain label are running
   - Disable a service temporarily, e.g. when developing it

### Useful links
 - Pseudo terminal for disabling buffering:
   https://github.com/kr/pty


## Apache License Version 2.0

   Copyright 2016 Daniel Josefsson

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
