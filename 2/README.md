# Example session for sending data via tcp
```
echo '49 00 00 30 39 00 00 00 65 49 00 00 30 3a 00 00 00 66 49 00 00 30 3b 00 00 00 64 49 00 00 a0 00 00 00 00 05 51 00 00 30 00 00 00 40 00 00 00 00 65 ' | xxd -r -p | nc localhost 8080
```
