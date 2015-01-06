/*

Author: Kyle Laplante

This library logs messages to Loglens using scribe.

Example:
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

*/


package loglens

import (
  "encoding/json"
  "fmt"
  "github.com/artyom/scribe"
  "github.com/artyom/thrift"
  "os"
  "strings"
  "time"
)

type LoglensClient struct {
  Category string
  Scribe   *scribe.ScribeClient
}

type LogSource struct {
  Message   string `json:"message"`            // required
  Tag       string `json:"tag,omitempty"`      // optional
  Type      string `json:"type,omitempty"`     // optional
  Username  string `json:"username,omitempty"` // optional
  logSourceGenerated  // hostname and timestamp are private fields
}

type Log struct {
  Index  string     `json:"index"`  // required
  Source *LogSource `json:"source"` // required
  Type   string     `json:"type"`   // required
  logId  // id is a private field
}

type logSourceGenerated struct {
  Hostname  string `json:"hostname"`   // private
  Timestamp string `json:"@timestamp"` // private
}

type logId struct {
  Id string `json:"id"` // private
}

func NewLoglensClient() *LoglensClient {
  return LoglensClientFactory("loglens", NewScribeClientFactory("localhost", "1463"))
}

func (l *LoglensClient) SimpleLog(logType, message, index string) (r scribe.ResultCode, err error) {
  log := &Log{Source: &LogSource{Message: message}, Index: index, Type: logType}
  return l.Log(log)
}

func LoglensClientFactory(category string, scribe *scribe.ScribeClient) *LoglensClient {
  return &LoglensClient{Category: category, Scribe: scribe}
}

func NewScribeClientFactory(host, port string) *scribe.ScribeClient {
  transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

  tsocket, err := thrift.NewTSocket(host + ":" + port)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  transport := transportFactory.GetTransport(tsocket)
  transport.Open()

  return scribe.NewScribeClientFactory(transport, thrift.NewTBinaryProtocolFactoryDefault())
}

func (l *LoglensClient) Log(log *Log) (r scribe.ResultCode, err error) {
  // create the timestamp so we can show it to STDOUT
  log.Source.Timestamp = strings.Replace(time.Now().UTC().Format(time.RFC3339), "Z", ".000", 1)
  // print as early as we can so it gets shown quickly if using "go client.Log(log)"
  fmt.Println(fmt.Sprintf("[%s %s] %s", log.Source.Timestamp, log.Type, log.Source.Message))

  // create private struct values
  log.Source.Hostname, _ = os.Hostname()
  log.Id = createUuid()

  // decode the json object to a byte array
  message, _ := json.Marshal(log)

  //fmt.Println(string(message))  // this shows the actual json string being sent to loglens
  return l.RawLog(string(message))
}

func (l *LoglensClient) RawLog(message string) (r scribe.ResultCode, err error) {
  return l.Scribe.Log([]*scribe.LogEntry{&scribe.LogEntry{Category: l.Category, Message: message}})
}

func (l *LoglensClient) Info(log *Log) (r scribe.ResultCode, err error) {
  log.Type = "INFO"
  return l.Log(log)
}

func (l *LoglensClient) Error(log *Log) (r scribe.ResultCode, err error) {
  log.Type = "ERROR"
  return l.Log(log)
}

func (l *LoglensClient) Warn(log *Log) (r scribe.ResultCode, err error) {
  log.Type = "WARN"
  return l.Log(log)
}

func createUuid() string {
  f, _ := os.Open("/dev/urandom")
  b := make([]byte, 16)
  f.Read(b)
  defer f.Close()
  return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
