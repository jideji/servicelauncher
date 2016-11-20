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
