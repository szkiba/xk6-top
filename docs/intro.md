---
author: Iván SZKIBA
date: YYYY-MM-DD
---

# x k 6 - t o p 

> *Terminal based metrics dashboard for k6*


              ▲
              │        /\
              │   /\  /  \
              │  /  \/    \
              │ /          \
              │/            \
              └───────────────►

[Repository](https://github.com/szkiba/xk6-top)

[Slides    ](https://ivan.szkiba.hu/xk6-top/intro)

<!--
Before the terminal-based dashboard demo, allow me to summarize what it is on a few slides.
I promise it will be a short and simple presentation.
We can say minimalistic.

Well, since the topic is a terminal-based dashboard, I thought the presentation should also be in the terminal.
As I said, minimalistic presentation.
-->

---

## Why make a terminal-based dashboard?

```
     _______________
    │.-------------.│
    ││             ││
    ││   W h y ?   ││              - Many k6 users love the terminal (me too)
    ││ ~~~~~~~~~~~ ││
    ││_____________││              - The terminal can be used for more than running CLI
    '------. .------'
           │ │    _│/              - Quick way to get an overview of the k6 test run
           │ │  ."   ".
           │ │ /(O)-(O)\           - TUI requires less resources than the browser
          /_)││   /     │
          │_)││  '-     │          - Easy to use even on a remote computer
          \_)│\ '.___.' /   │\/│_
           │ │ \  \_/  /   _│  '/    - Doesn't require opening a TCP port
           │_│\ '.___.'    \ ) /
           \   \_/\__/\__   │==│     - Doesn't require preparation (e.g. database)
            \    \ /\ /\ `\ │  │
             \    \\//     \│  │
              `\   /\   │  /   │
                ;  ││   │\____/
                │  ││   │
```
<!--
Now that the web-dashboard has been included as an experimental module in k6, the question arises why a terminal-based dashboard is needed.
Does it make sense, what is it good for?

First, many k6 users (myself included) like to work in a terminal.
The terminal can be used for much more than running CLI programs.
This is what the text-based user interface or terminal-based user interface topic is all about.
Today's terminals, more specifically terminal emulators, can handle 256 or more colors, have unicode support and so on.

The status of the k6 test run can be quickly reviewed in the terminal.
The terminal-based dashboard uses fewer resources (CPU, memory) than a web dashboard running in a browser.
This is an important advantage if the user runs the k6 test and uses the dashboard on the same computer.

It is easy to use the dashboard in the terminal even in the case of a k6 test run on a remote computer.
It is not necessary to open or forward a TCP port.
The terminal emulator software securely transmits data to your screen from the remote computer.
-->

---
## What exactly is xk6-top?

```
     _______________
    │.-------------.│
    ││             ││
    ││  W h a t ?  ││              - Terminal-based metrics dashboard for k6
    ││ ~~~~~~~~~~~ ││
    ││_____________││                  - Quick overview of the k6 test run
    '------. .------'
           │ │    _│/                  - Current status of thresholds
           │ │  ."   ".
           │ │ /(O)-(O)\               - Table display of metrics
          /_)││   /     │
          │_)││  '-     │              - Chart display of relevant metrics
          \_)│\ '.___.' /   │\/│_
           │ │ \  \_/  /   _│  '/      - All values are updated dynamically
           │_│\ '.___.'    \ ) /
           \   \_/\__/\__   │==│   - k6 output extension (k6 run -o top script.js)
            \    \ /\ /\ `\ │  │
             \    \\//     \│  │   - k6 launcher (k6top run script.js)
              `\   /\   │  /   │
                ;  ││   │\____/    - Web dashboard viewer (k6top [-u URL])
                │  ││   │
```
<!--
So what exactly is the xk6-top?

Well, it's basically a k6 metrics dashboard that works in the terminal.

It enables a quick overview of the status of the k6 test run.

Displays the current status of the thresholds in detail.
The individual threshold expressions are marked with green (passed) or red (failed) color according to the success of the evaluation.
From this table, you can immediately see which threshold expression failed.

Displays the k6 metrics in tabular form.

Draws simple graphs of relevant metrics. Of course, within the limits given by the terminal.

All displayed values are dynamically updated during the test run.
The dashboard can be used in three ways.

It can be used as a k6 output extension. In this case, the name "top" must be specified with the output flag in the k6 run command.

It can be used as a k6 launcher. k6top is a standalone program independent of k6. The "run" subcommand can be used to start k6. In this case, it is not necessary to build a custom k6, since the web-dashboard has been part of k6 since version v0.49.0.

It can be used as a dashboard viewer. The k6top command can be given the address of the k6 web dashboard. By default, it uses the local k6 web dashboard, but it is suitable for displaying any remote k6 metrics, if the web dashboard feature is enabled in k6.
-->

---

## How It Works?

```
     _______________
    │.-------------.│
    ││             ││
    ││   H o w ?   ││              - Utilizes web dashboard's SSE stream
    ││ ~~~~~~~~~~~ ││
    ││_____________││                - Acts as a web dashboard client
    '------. .------'
           │ │    _│/              - Parses and evaluates thresholds expressions
           │ │  ."   ".
           │ │ /(O)-(O)\             - Thresholds status is updated every 10 seconds
          /_)││   /     │
          │_)││  '-     │          - Uses bubbletea TUI framework
          \_)│\ '.___.' /   │\/│_
           │ │ \  \_/  /   _│  '/    - Dashboard events become messages
           │_│\ '.___.'    \ ) /
           \   \_/\__/\__   │==│
            \    \ /\ /\ `\ │  │
             \    \\//     \│  │
              `\   /\   │  /   │
                ;  ││   │\____/
                │  ││   │
```
<!--
The terminal-based dashboard uses the same SSE stream as a data source as the web dashboard.
To be precise, the SSE stream is generated by the k6 extension part of the web dashboard and used by the web dashboard UI part and the terminal dashboard UI.
In this sense, the terminal-based dashboard can also be considered the client of the web dashboard running in the terminal.

One of the most convenient features is the parsing and continuous evaluation of threshold expressions based on current metrics.
The evaluation is done using a simple expression evaluation library and is compatible with the k6 threshold expression evaluator (at least I hope so).

The user interface was created using the bubbletea framework.
This is a simple, easy-to-understand library.
Nowadays, this is one of the most popular text-based user interface libraries for the go language.
The operation of bubbletea is message-based, so it can be well matched to the event-based dashboard data stream.
Dashboard events become bubbletea messages and the screen is updated based on these messages.

The k6 binary size increase caused by bubbletea is acceptable. I didn't measure exactly how much, but the total size increase in the case of xk6-top is approximately 1MB, which also includes other libraries.
-->
