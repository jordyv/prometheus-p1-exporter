# Protocol description #

## DSMR 5.0 ##

```
1   /Ene5\XS210 ESMR 5.0

2   1-3:0.2.8(50)
3   0-0:1.0.0(181009214805S)
4   0-0:96.1.1(xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
5   1-0:1.8.1(001179.186*kWh)
6   1-0:1.8.2(001225.590*kWh)
7   1-0:2.8.1(000000.016*kWh)
8   1-0:2.8.2(000000.000*kWh)
9   0-0:96.14.0(0002)
10  1-0:1.7.0(00.200*kW)
11  1-0:2.7.0(00.000*kW)
12  0-0:96.7.21(00057)
13  0-0:96.7.9(00002)
14  1-0:99.97.0(1)(0-0:96.7.19)(170829233732S)(0000001803*s)
15  1-0:32.32.0(00002)
16  1-0:32.36.0(00000)
17  0-0:96.13.0()
    1-0:32.7.0(227.0*V)
    1-0:31.7.0(001*A)
    1-0:21.7.0(00.200*kW)
    1-0:22.7.0(00.000*kW)
18  0-1:24.1.0(003)
19  0-1:96.1.0(xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
20  0-1:24.2.1(181009214500S)(01019.003*m3)
    !6611
```
```
1   Header
2   Version
3   Timestamp
4   Equipment ID
5   Electricity delivered to client tariff 1 (per 0.001kWh)
6   Electricity delivered to tariff 2 (per 0.001kWh)
7   Electricity delivered by client tariff 1 (per 0.001kWh)
8   Electricity delivered by tariff 2 (per 0.001kWh)
9   Tariff indicator
10  Actual electricity power delivered in 1 W resolution
11  Actual electricity power received in 1 W resolution
12  Number of power failures in any phase
13  Number of long power failures in any phase
14  Power failure event log (long power failures)
        end of failure - duration in seconds
15  Number of voltage sags in phase L1
16  Number of voltage swells in phase L1
17  Text message
18  Device type
19  Gas meter equipment ID
20  Last hourly value of gas delivered to client in m3 with timestamp
```

[source](https://www.netbeheernederland.nl/_upload/Files/Slimme_meter_15_32ffe3cc38.pdf)
