# chatcontrol

Build: install golang, go to project directory and run this:
```
go install
```

Integration:
```
from ably import AblyRest

ably = AblyRest('get ur key')
channel = ably.channels.get('cvm')

if ws_message[0] == "chat" and ws_message[2].startswith(".run "):
  command = ws_message[2][5:]
  print("lol", command)

  await channel.publish("shell", command)

```
