arduino-cli compile --fqbn arduino:avr:uno --build-path ./sketch/build/ ./sketch/sketch.ino

avrdude -v -patmega328p -carduino -PCOM5 -b115200 -D -Uflash:w:sketch/build/sketch.ino.hex:i