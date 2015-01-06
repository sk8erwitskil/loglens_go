### Loglens ###

This library logs messages to Loglens using scribe.
 
Example:
```
  // First create a Loglens Client
  client := loglens.NewLoglensClient()
 
  // Some string to send to loglens
  message := "Here is a Loglens Message"
 
  // You can be very explicit by creating the LogSource and Log structs manually
  source := &loglens.LogSource{Message: message, Username: "peacock", Tag: "golang", Type: "INFO"}
  log := &loglens.Log{Index: "peacock", Source: source, Type: "INFO"}
  result, err := client.Log(log)
  fmt.Println(result, err)
 
  // You can use the built-in .Info, .Warn and .Error methods from LoglensClient
  // to skip having to manually define the Type.
  source = &loglens.LogSource{Message: message, Username: "peacock"}
  log = &loglens.Log{Index: "peacock", Source: source}
  result, err = client.Error(log)
  fmt.Println(result, err)
 
  // You can skip having to create the LogSource and Log structs by
  // using the .SimpleLog method from LoglensClient.
  result, err = client.SimpleLog("WARN", message, "peacock")
  fmt.Println(result, err)
 ```
