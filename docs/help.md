### Navigation

The screen is divided into so-called tabs. You can switch between tabs using the navigation bar at the top of the screen. 

Key                  | Function
---------------------|---------
`Esc`, `Ctrl+c`, `q` | Quit
`Right`, `Tab`       | Switch to the next tab
`Left`, `Shift+Tab`  | Switch to the previous tab
`Shift+Right`        | Move to next chart
`Shift+Left`         | Move to previous chart
`Down`,`PgDown`      | Move down
`Up`, `PgUp`         | Move up
`+`, `Shift+Down`    | Expand, show more details
`-`, `Shift+Up`      | Collapse, show less details

There are also charts on some tabs. Due to the limitations of the terminal, one chart is displayed at a time. You can switch between charts with the `Shift+Right` and `Shift+Left` keys.

Certain tabs (on which data is not available) may be disabled. Their names appear in italics.

### Overview tab

The *Overview* tab provides an overview of the test run. Here you can find the most important parameters such as start time, elapsed time, remaining time (in the case of a running test) and the colored result of the test run (in the case of a completed test).

The most important element of the *Overview* tab is the table containing the current state of the thresholds in detail. The individual threshold expressions are marked with green (passed) or red (failed) color according to the success of the evaluation. From this table, you can immediately see which threshold expression failed.

### Metrics tables

The aggregated metric values are available in tabular form on the tabs named according to the type of metric (*Trends*, *Counters*, *Rates*, *Gauges*).

By default, the table contains the tags expanded. Rows containing tags can be collapsed with the `-` key or the `Shift+Up` key combination. They can be expanded at any time later with the `+` key or the `Shift+Down` key combination.

### Metrics charts

The most important metrics are also available in the form of a time chart on the tabs with names corresponding to the protocol (*HTTP*, *gRPC*, *WS*, *Browser*).

Due to the limitations of the terminal, one chart can be seen on these tabs at a time. A second-level navigation enables the choice of a chart within a tab. You can switch to the next/previous chart with the key combinations `Shift+Right` and `Shift+Left`.

The chart displaying the *trend* type metric contains the percentile values in addition to the average value. Percetile series can be hidden individually by repeatedly using the `-` key or the `Shift+Up` key combination. They can be displayed at any type later with repeatedly using the `+` key or the `Shift+Down` key combination.
