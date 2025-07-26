String input = "";
const int MAX_TOKENS = 10;
const int speakerPin = 9;

const int LEDCUBE_L1 = 6;
const int LEDCUBE_L2 = 7;

const int LEDCUBE_C1 = 2;
const int LEDCUBE_C2 = 3;
const int LEDCUBE_C3 = 4;
const int LEDCUBE_C4 = 5;

const int CUBE_POINT[6] = {LEDCUBE_L2, LEDCUBE_L1, LEDCUBE_C4, LEDCUBE_C3, LEDCUBE_C2, LEDCUBE_C1};

void setup() {
  Serial.begin(9600);
  pinMode(LED_BUILTIN, OUTPUT);
  pinMode(LEDCUBE_C1, OUTPUT);
  pinMode(LEDCUBE_C2, OUTPUT);
  pinMode(LEDCUBE_C3, OUTPUT);
  pinMode(LEDCUBE_C4, OUTPUT);
  pinMode(LEDCUBE_L1, OUTPUT);
  pinMode(LEDCUBE_L2, OUTPUT);
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
  } else if (args[0] == "beep") {
    if (count < 3) return invaidargs();
    beep(args[1],args[2]);
  } else if (args[0] == "cube") {
    if (count < 2) return invaidargs();
    cube(args[1]);
  } else {
    invaidcmd();
  }
}

void cube(String state) {
  for (int i = 0; i < 6; i++) {
    if (state.charAt(i) == '0') {
      digitalWrite(CUBE_POINT[i], LOW);
    } else if (state.charAt(i) == '1') {
      digitalWrite(CUBE_POINT[i], HIGH);
    }
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

void beep(String frequency, String duration) {
  tone(speakerPin, frequency.toInt(), duration.toInt());
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