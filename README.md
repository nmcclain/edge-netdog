# edge-netdog

`edge-netdog` performs emergency remediation (such as a reboot) for edge devices when network connectivity has failed.


It's useful on edge devices, like Raspberry Pi's, that need a last resort when unstable networks or unfriendly WiFi connections are a problem. 

WARNING: This tool is intended for **edge** devices, and is **not appropriate for typical servers.**

Read [more about the motivation here](https://medium.com/@nedmcclain/edge-networkings-last-resort-890b536ab960).

## Usage

Download the [latest version here](https://github.com/nmcclain/edge-netdog/releases).

### Required/Minimal Configuration File

* The `target_url` and `target_match` fields are required in the `global` section.

```
---
global:
    target_url: https://example.com
    target_match: Example Domain

actions:
    - echo network down
```

### Full Configuration File
The remediation actions shown here may be appropriate for a Raspberry Pi on WiFi:
```
---
global:
    debug: true
    monitor_interval: 30s
    target_attempts: 4
    action_delay: 30s
    target_url: https://example.com
    target_match: Example Domain

actions:
    - sudo wpa_cli -i wlan0 reconfigure
    - sudo dhclient -v
    - sudo service networking restart
    - sudo reboot
```

* You can test failures by changing the `target_match` to a random value.
* For production, use a domain you control rather than example.com.

## Contributing
**Something bugging you?** Please open an [Issue](https://github.com/nmcclain/edge-netdog/issues) or [Pull Request](https://github.com/nmcclain/edge-netdog/pulls) - we're here to help!

**New Feature Ideas?** Please open a [Pull Request](https://github.com/nmcclain/edge-netdog/pulls).
 
**All Humans Are Equal In This Project And Will Be Treated With Respect.**
