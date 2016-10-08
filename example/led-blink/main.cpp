#include "app.h"


void setup() {
    static int i = 0;

    Timer.set_interval(CONFIG_BLINK_INTERVAL, []() {
        if (i % 2 == 0)
            MainLED.on();
        else
            MainLED.off();

        i++;
    });
}
