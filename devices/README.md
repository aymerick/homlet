homlet devices
==============

DIY devices for [homlet](https://github.com/aymerick/homlet).


Jeenode
=======

Jeenode design by JC Wippler: <http://jeelabs.net/projects/hardware/wiki/JeeNode>

References
----------

  - <http://jeelabs.net/projects/hardware/wiki/JeeNode>
  - <http://jeelabs.org/2011/12/13/developing-a-low-power-sketch/>
  - <http://jeelabs.org/2011/06/09/rf12-packet-format-and-design/>
  - <http://jeelabs.org/2011/01/14/nodes-addresses-and-interference/>
  - <http://jeelabs.org/2010/12/07/binary-packet-decoding/>
  - <http://jeelabs.org/2010/12/08/binary-packet-decoding-–-part-2/>
  - <http://jeelabs.org/2013/09/05/decoding-bit-fields/>
  - <http://jeelabs.org/2013/09/06/decoding-bit-fields-part-2/>

TinyTX
======

TinyTX design by Nathan Chantrell: <http://nathan.chantrell.net/tinytx-wireless-sensor>

Setup
-----

Arduino IDE Setup (on my Mac):

  - Install [OneWire lib](http://www.pjrc.com/teensy/arduino_libraries/OneWire.zip) into `~/Documents/Arduino/libraries/OneWire/`
  - [Fix OneWire lib](http://arduino.cc/forum/index.php/topic,91491.msg687523.html#msg687523) - Open `OneWire.h` and below the line:
```
      #include “Arduino.h”       // for delayMicroseconds, digitalPinToBitMask, etc
```
add:
```
      #include “pins_arduino.h”  // for digitalPinToBitMask, etc
```
  - Install [DallasTemperature lib](http://download.milesburton.com/Arduino/MaximTemperature/DallasTemperature_LATEST.zip) into `~/Documents/Arduino/libraries/DallasTemperature/`
  - Install [Jeelib](https://github.com/jcw/jeelib) into `~/Documents/Arduino/libraries/jeelib/`

Programmer setup:

  - Use an ISP USB programmer
  - Setup an ISP adapter: <http://www.evilmadscientist.com/2007/using-avr-microcontrollers-minimalist-target-boards/>

Program ATTiny84:

  - Setup ATTiny support in Arduino thanks to: <http://code.google.com/p/arduino-tiny>
  - Select `Tools/Board` > `ATtiny84 @ 8Mhz (internal oscillator; BOD disabled)`
  - Select `Tools/Programmer/USBasp`
  - Select `Tools/Burn Bootloader`. Note that this isn’t actually burning a bootloader to the ATtiny (it doesn’t use one), it is just using this function to set the AVR fuses to configure the oscillator at 8MHz.
  - Upload a `TinyTX_*` sketch

Build circuit:

  - You can try on a stripboard with: <http://nathan.chantrell.net/downloads/arduino/tinytx/tinytx_stripboard_ds18b20.png>
  - You can too order PCBs from SeeedStudio as explained here: <http://nathan.chantrell.net/tinytx-wireless-sensor/#div-comment-190398>

Issues:

  - If the sketch is stuck during call to `rf12_sendWait` then Burn Bootloader to fix the issue

References
----------

  - <https://github.com/nathanchantrell/TinyTX/blob/master/TinyTX_DHT22/TinyTX_DHT22.ino>
  - <https://github.com/mharizanov/TinySensor/blob/master/Funky_DHT22/Funky_DHT22.ino>
