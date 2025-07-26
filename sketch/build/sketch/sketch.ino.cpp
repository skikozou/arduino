#include <Arduino.h>
#line 1 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
String input = "";
const int MAX_TOKENS = 10;
int speakerPin = 8;

#line 5 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void setup();
#line 11 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void loop();
#line 24 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void router(String cmd);
#line 41 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void led(String state);
#line 52 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void invaidcmd();
#line 56 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void invaidargs();
#line 60 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
int split(String str, char delimiter, String tokens[]);
#line 5 "C:\\Users\\skiko\\go\\arduino\\sketch\\sketch.ino"
void setup() {
  Serial.begin(9600);
  pinMode(LED_BUILTIN, OUTPUT);
  tone(speakerPin, 440, 1000);
}

void loop() {
  while (Serial.available()) {
    char c = Serial.read();
    if (c == '\n') {
      router(input);
      
      input = "";
    } else {
      input += c;
    }
  }
}

void router(String cmd) {
  cmd.trim();
  String args[10];
  int count = split(cmd, ' ', args);

  if (count < 1) {
    return invaidargs();
  }

  if (args[0] == "led") {
    if (count < 2) return invaidargs();
    led(args[1]);
  } else {
    invaidcmd();
  }
}

void led(String state) {
  state.trim();
  if (state == "high") {
    digitalWrite(LED_BUILTIN, HIGH);
  } else if (state == "low") {
    digitalWrite(LED_BUILTIN, LOW);
  } else {
    invaidargs();
  }
}

void invaidcmd() {
  Serial.println("invaid command");
}

void invaidargs() {
  Serial.println("invaid args");
}

int split(String str, char delimiter, String tokens[]) {
  int index = 0;
  int start = 0;
  int end = str.indexOf(delimiter);

  while (end != -1 && index < MAX_TOKENS) {
    tokens[index++] = str.substring(start, end);
    start = end + 1;
    end = str.indexOf(delimiter, start);
  }

  if (index < MAX_TOKENS) {
    tokens[index++] = str.substring(start);
  }

  return index;
}
