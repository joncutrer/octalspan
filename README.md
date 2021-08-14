# octalspan

A modern open source syslog server written in Go

:warning: DISCLAIMER: Do not use this software for production environments. :warning:

## About the project

This project is currently an experiment in it's early stages.  The golang library `gopkg.in/mcuadros/go-syslog` currently does all low-level handling of syslog messages.



## Installation Steps

* Run installer **octalspan-0.#.#.exe**
* Open administrative command prompt
* Change to application directory
  ```
  cd "c:\Program Files (x86)\octalspan"
  ````
* Install windows service
  ```
  octalspan.svc.exe install
  ```
* Start service
  ```
  octalspan.svc.exe start
  ```

## Post Installation Notes

* You must open ports in Windows Firewall before the service can receive syslog messages. 
  
  * 514/udp
  * 601/tcp 


* IP bindings and listening ports can be customized by editing `octalspan.yml`

   


* Received Syslog messages are currently written to `c:\Program Files (x86)\octalspan\log`

