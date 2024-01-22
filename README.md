# findmyticker

### build

```go build -o ./build/FindMy\ ticker.app/Contents/MacOS/findmyticker .```

### run

Allow run by adding `findmyticker` app in `Settings->Privacy&Security->Full disk access` by `+`

and then 

`xattr -d com.apple.quarantine /Applications/FindMy\ ticker.app`

### How it works

While FindMy is running, it stores the information of your devices in temporary files `~/Library/Caches/com.apple.findmy.fmipcore/Items.data` and `~/Library/Caches/com.apple.findmy.fmipcore/Devices.data` in JSON format. The app periodically reads these files and stores a new record for each device if there is an update.