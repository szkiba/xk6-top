Display k6 metrics on a terminal-based dashboard

`k6top` connects to a k6 process (even a remote one) and displays metrics on a terminal-based dashboard.

It utilizes the SSE stream of the k6 web dashboard to display the aggregated metric values. The aggregation period is 10s by default. The data on the dashboard is updated at the end of each period. This period can be modified by configuring the k6 web dashboard.

The address of the k6 web dashboard can be specified with the `--url` flag. Thus, the metrics of local k6 running on a non-default port or even k6 running on a remote computer can be displayed.
