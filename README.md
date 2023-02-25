# keysnogger

![GitHub](https://img.shields.io/github/license/jibstack64/keysnogger) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jibstack64/keysnogger)

A cross-platform keylogger.

![Preview](https://github.com/jibstack64/keysnogger/blob/master/preview.png)

One minor issue with the keylogger is the tracking of backspace, tab and enter key presses. They are sent to the server as `<backspace>`, `<tab>` and `<enter>` and then decoded when being write to the log file for that specific client. If someone knew that the keylogger was running (very unlikely), they could type these "codes" so to speak, and erase portions of the logs.