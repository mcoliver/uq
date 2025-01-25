# uq
In the spirit of jq but for parsing url's to key value pairs and converting url encoding to ascii

## Dependencies
https://github.com/fatih/color

## Build instructions
```
go mod init uq
go mod tidy
CGO_ENABLED=0 go build -ldflags="-s -w" -o uq main.go
```

## Run it
```
./uq "https://github.com/page?redirect=https%3A%2F%2Fmicrosoft.com%2Fredirect%3Ftarget%3Dhttps%253A%252F%252Fgoogle.com%252Fpath%253Fkey1%253Dvalue1%2526key2%253Dvalue2%26source%3Dapp&user_id=12345&session=abc123%26extra_param%3Dtrue"
"
```
```
{
    "base_url":  "https://github.com/page",
    "query":  {
      "redirect":  {
        "base_url":  "https://microsoft.com/redirect",
        "query":  {
          "source":  "app",
          "target":  {
            "base_url":  "https://google.com/path",
            "query":  {
              "key1":  "value1",
              "key2":  "value2"
          }
        }
      }
    },
      "session":  "abc123\u0026extra_param=true",
      "user_id":  "12345"
  }
}
```
