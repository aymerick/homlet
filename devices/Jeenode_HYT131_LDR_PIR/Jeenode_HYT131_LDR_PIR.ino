//
// Homlet Device
//
// Jeenode with sensors:
//  - HYT131 for temperature and humidity
//  - LDR for light
//  - PIR for motion
//
// References:
//  - Room board: http://jeelabs.com/products/room-board
//  - HYT131: https://jeelabs.org/2012/06/30/new-hyt131-sensor/index.html
//  - Original sketch: https://github.com/jcw/jeelib/blob/master/examples/RF12/roomNode/roomNode.ino
//

#include <JeeLib.h>      // https://github.com/jcw/jeelib
#include <PortsSHT11.h>
#include <avr/sleep.h>
#include <util/atomic.h>

// has to be defined because we're using the watchdog for low-power waiting
ISR(WDT_vect) { Sleepy::watchdogEvent(); }


#define DEBUG 0
#define NOOP  0

// Device kind
#define DEVICE_KIND 1

// defined if HYT131 is connected to a port
#define HYT131_PORT 1

// defined if LDR is connected to a port's AIO pin
#define LDR_PORT 4

// defined if PIR is connected to a port's DIO pin
#define PIR_PORT 4

// set to one to pull-up the PIR input pin
#define PIR_PULLUP 1

// 0 or 1, to match PIR reporting high or low
#define PIR_INVERTED 0

// how often to measure, in tenths of seconds
#define MEASURE_PERIOD 600

// number of milliseconds to wait for an ack
#define ACK_TIME 10

// how soon to retry if ACK didn't come in
#define ACK_RETRY_PERIOD 10

// maximum number of times to retry
#define ACK_RETRY_LIMIT 5

// report every N measurement cycles
#define REPORT_EVERY 5

// set the sync mode to 2 if the fuses are still the Arduino default
// mode 3 (full powerdown) can only be used with 258 CK startup fuses
#define RADIO_SYNC_MODE 2


// The scheduler makes it easy to perform various tasks at various times
enum { MEASURE, REPORT, TASK_END };

static word schedbuf[TASK_END];
Scheduler scheduler (schedbuf, TASK_END);

// count up until next report, i.e. packet send
static byte reportCount;

static byte myDeviceID;

// serialized payload
struct {
  byte kind     :7;  // Device kind
  byte reserved :1;  // Reserved for future use. Must be zero.
  // data
  int  temp     :10; // Temperature: -500..+500 (tenths)
  byte humi     :7;  // Humidity: 0..100
  byte light    :8;  // Light sensor: 0..255
  byte motion   :1;  // Motion detector: 0..1
  byte lowbat   :1;  // Supply voltage dropped under 3.1V: 0..1
} payload;

// Interface to a Passive Infrared motion sensor.
class PIR : public Port {
  volatile byte value, changed;

  public:
    PIR (byte portnum)
      : Port (portnum), value (0), changed (0) {}

    // this code is called from the pin-change interrupt handler
    void poll() {
      // see http://talk.jeelabs.net/topic/811#post-4734 for PIR_INVERTED
      byte pin = digiRead() ^ PIR_INVERTED;

      // if the pin changed then set the changed flag to report it
      if (pin != state()) {
        changed = 1;
      }

      value = pin;
    }

    // state is true if curr value is still on or if it was on recently
    byte state() const {
      byte f = value;
      return f;
    }

    // return true if PIR state changed
    byte stateChanged() {
      byte f = changed;
      changed = 0;
      return f;
  }
};

// sensors
PIR sensorPIR (PIR_PORT);
PortI2C hyti2cport (HYT131_PORT); // FIXME rename
HYT131 sensorHYT131 (hyti2cport); // FIXME rename
Port sensorLDR (LDR_PORT);

// the PIR signal comes in via a pin-change interrupt
ISR(PCINT2_vect) { sensorPIR.poll(); }


//
// Helpers
//

static void serialFlush () {
#if ARDUINO >= 100
  Serial.flush();
#endif

  // make sure tx buf is empty before going back to sleep
  delay(2);
}

// wait a few milliseconds for proper ACK to me, return true if indeed received
static byte waitForAck() {
  MilliTimer ackTimer;

  while (!ackTimer.poll(ACK_TIME)) {
    // see http://talk.jeelabs.net/topic/811#post-4712
    if (rf12_recvDone() && (rf12_crc == 0) && (rf12_hdr == (RF12_HDR_DST | RF12_HDR_CTL | myDeviceID)))
      return 1;

    set_sleep_mode(SLEEP_MODE_IDLE);
    sleep_mode();
  }

  return 0;
}

// send payload and wait for ack
static void sendPayload() {
  for (byte i = 0; i < ACK_RETRY_LIMIT; i++) {
    // power up RF
    rf12_sleep(RF12_WAKEUP);

    // send payload
    rf12_sendNow(RF12_HDR_ACK, &payload, sizeof payload);
    rf12_sendWait(RADIO_SYNC_MODE);

    // wait for ack
    byte acked = waitForAck();

    // power down RF
    rf12_sleep(RF12_SLEEP);

    if (acked) {
#if DEBUG
      Serial.print(" ack");
      Serial.println((int) i);
      serialFlush();
#endif

      // reschedule measurements
      scheduler.timer(MEASURE, MEASURE_PERIOD);

      return;
    }

    // if no ack received wait and try again
    delay(ACK_RETRY_PERIOD * 100);
  }

#if DEBUG
  Serial.println(" no ack!");
  serialFlush();
#endif

  // reschedule measurements
  scheduler.timer(MEASURE, MEASURE_PERIOD);
}

// send payload and optionally report on serial port
static void doReport() {
#if DEBUG
  Serial.print((int) payload.light);
  Serial.print(' ');
  Serial.print((int) payload.motion);
  Serial.print(' ');
  Serial.print((int) payload.humi);
  Serial.print(' ');
  Serial.print((int) payload.temp);
  Serial.print(' ');
  Serial.print((int) payload.lowbat);
  Serial.println();
  serialFlush();
#endif

#if !NOOP
  sendPayload();
#endif
}


//
// Sensors
//

// read LDR data
void readLDR() {
  // enable AIO pull-up
  sensorLDR.digiWrite2(1);

  payload.light = ~ sensorLDR.anaRead() >> 2;

  // disable pull-up
  sensorLDR.digiWrite2(0);
}

// read Battery status
void readLowBat() {
  payload.lowbat = rf12_lowbat();
}

// read temperature and humidity
void readHYT131() {
  int humi, temp;
  sensorHYT131.reading(temp, humi);
  payload.humi = humi/10;
  payload.temp = temp;
}

// readout all the sensors and other values
static void readSensors() {
  readLDR();
  readLowBat();
  readHYT131();

  payload.motion = sensorPIR.state();
}


//
// Main
//

void setup() {
#if DEBUG
  Serial.begin(57600);
  myDeviceID = rf12_config();
  serialFlush();
#else
  // don't report info on the serial port
  myDeviceID = rf12_config(0);
#endif

  // power down
  rf12_sleep(RF12_SLEEP);

  sensorPIR.digiWrite(PIR_PULLUP);

#ifdef PCMSK2
  bitSet(PCMSK2, PIR_PORT + 3);
  bitSet(PCICR, PCIE2);
#endif

  // init payload
  payload.reserved = 0;
  payload.kind = DEVICE_KIND;

  // report right away for easy debugging
  reportCount = REPORT_EVERY;

  // start the measurement loop going
  scheduler.timer(MEASURE, 0);
}

void loop() {
#if DEBUG
  Serial.print('.');
  serialFlush();
#endif

  if (sensorPIR.stateChanged()) {
    payload.motion = sensorPIR.state();

#if DEBUG
    Serial.print("PIR: ");
    Serial.print((int) payload.motion);
    Serial.print(" ");
    serialFlush();
#endif

    // report
    doReport();
  }

  // wait for events
  switch (scheduler.pollWaiting()) {
  case MEASURE:
    // reschedule measurements
    scheduler.timer(MEASURE, MEASURE_PERIOD);

    // read sensors
    readSensors();

    // every so often, a report needs to be sent out
    if (++reportCount >= REPORT_EVERY) {
      reportCount = 0;

      // schedule report
      scheduler.timer(REPORT, 0);
    }
    break;

  case REPORT:
    // report
    doReport();
    break;
  }
}
