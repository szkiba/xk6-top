Run k6 test and display metrics on terminal-based dashboard

The `k6top run` command starts the `k6 run` command with the specified arguments and then displays the metrics dashboard in the terminal. A `k6` with at least version `v0.49.0` must be in the command search path.

In the launched k6, the web dashboard feature will be enabled. This is necessary because the SSE stream of the web dashboard is also used by the terminal-based dashboard as a data source.

If the environment variables `K6_WEB_DASHBOARD_HOST` and `K6_WEB_DASHBOARD_PORT` are set, their values are taken into account by the command.

The `k6top run` command has no flags, all arguments are passed to the `k6 run` command without interpretation or changes.
