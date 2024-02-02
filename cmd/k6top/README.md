<!-- #region cli -->
## k6top

Terminal based metrics dashboard viewer for k6

### Synopsis

Display k6 metrics on a terminal-based dashboard

`k6top` connects to a k6 process (even a remote one) and displays metrics on a terminal-based dashboard.

It utilizes the SSE stream of the k6 web dashboard to display the aggregated metric values. The aggregation period is 10s by default. The data on the dashboard is updated at the end of each period. This period can be modified by configuring the k6 web dashboard.

The address of the k6 web dashboard can be specified with the `--url` flag. Thus, the metrics of local k6 running on a non-default port or even k6 running on a remote computer can be displayed.


```
k6top [flags]
```

### Flags

```
  -h, --help         help for k6top
  -u, --url string   k6 web dashboard URL (default "http://127.0.0.1:5665")
```

### SEE ALSO

* [k6top run](#k6top-run)	 - k6 test runner and terminal-based metrics dashboard viewer

---
## k6top run

k6 test runner and terminal-based metrics dashboard viewer

### Synopsis

Run k6 test and display metrics on terminal-based dashboard

The `k6top run` command starts the `k6 run` command with the specified arguments and then displays the metrics dashboard in the terminal. A `k6` with at least version `v0.49.0` must be in the command search path.

In the launched k6, the web dashboard feature will be enabled. This is necessary because the SSE stream of the web dashboard is also used by the terminal-based dashboard as a data source.

If the environment variables `K6_WEB_DASHBOARD_HOST` and `K6_WEB_DASHBOARD_PORT` are set, their values are taken into account by the command.

The `k6top run` command has no flags, all arguments are passed to the `k6 run` command without interpretation or changes.


```
k6top run [flags]
```

### Flags

```
  -h, --help   help for run
```

### SEE ALSO

* [k6top](#k6top)	 - Terminal based metrics dashboard viewer for k6

<!-- #endregion cli -->
