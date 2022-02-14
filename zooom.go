package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "os"
  "os/exec"
  "runtime"
  "strconv"
  "strings"
  "time"
)

type Buffer struct {
  Start string `json:"start"`
  End   string `json:"end"`
}

type Config struct {
  Sources []string `json:"sources"`
  Buffer  Buffer   `json:"buffer"`
}

type Meeting struct {
  Name  string   `json:"name"`
  Days  []string `json:"days"`
  Start string   `json:"start"`
  End   string   `json:"end"`
  Metno string   `json:"metno"`
  Paswd string   `json:"paswd"`
}

type Now struct {
  fst   string
  wkd   string
  start time.Time
  end   time.Time
}

func read_config(loct string) (config Config) {
  if strings.HasPrefix(loct, "http://") || strings.HasPrefix(loct, "https://") {
    resp, _ := http.Get(loct)
    decoder := json.NewDecoder(resp.Body)
    decoder.Decode(&config)
  } else {
    file, _ := os.Open(loct)
    defer file.Close()
    decoder := json.NewDecoder(file)
    decoder.Decode(&config)
  }
  return
}

func read_meetings(loct string) (meets []Meeting) {
  if strings.HasPrefix(loct, "http://") || strings.HasPrefix(loct, "https://") {
    resp, _ := http.Get(loct)
    decoder := json.NewDecoder(resp.Body)
    decoder.Decode(&meets)
  } else {
    file, _ := os.Open(loct)
    defer file.Close()
    decoder := json.NewDecoder(file)
    decoder.Decode(&meets)
  }
  return
}

func bufferize(tm time.Time, bst string, add bool) (ntm time.Time) {
  buff := []int{0, 0}
  for i, s := range strings.Split(bst, ":") {
    val, _ := strconv.ParseInt(s, 10, 0)
    buff[i] = int(val)
  }
  if add {
    ntm = tm.Add(time.Hour*time.Duration(buff[0]) + time.Minute*time.Duration(buff[1]))
  } else {
    ntm = tm.Add(time.Hour*time.Duration(-buff[0]) + time.Minute*time.Duration(-buff[1]))
  }
  return
}

func find_meets(conf Config) (meets_now []Meeting) {
  nt := time.Now()
  now := Now{
    fst:   nt.Format(time.RFC3339),
    wkd:   nt.Weekday().String()[:3],
    start: bufferize(nt, conf.Buffer.Start, true),
    end:   bufferize(nt, conf.Buffer.End, false),
  }
  var start, end time.Time

  for _, source := range conf.Sources {
    meets := read_meetings(source)
  meet:
    for _, meet := range meets {
      for _, day := range meet.Days {
        if day == now.wkd {
          start, _ = time.Parse(time.RFC3339, now.fst[:11]+meet.Start+":00"+now.fst[19:])
          end, _ = time.Parse(time.RFC3339, now.fst[:11]+meet.End+":00"+now.fst[19:])
          if !now.start.Before(start) && now.end.Before(end) {
            meets_now = append(meets_now, meet)
          } else {
            continue meet
          }
        }
      }
    }
  }
  return
}

func choose_meet(meets []Meeting) (meet Meeting) {
  switch len(meets) {
  case 0:
    return
  case 1:
    meet = meets[0]
  default:
    fmt.Println(len(meets), "Meetings found:")
    for i, meet := range meets {
      fmt.Printf(" [%d] %s\n", i+1, meet.Name)
    }
    fmt.Println()
    var ch int
    for i := 0; i < 3; i++ {
      fmt.Print("Choose Meeting: ")
      fmt.Scanf("%d\n", &ch)
      if 1 <= ch && int(ch) <= len(meets) {
        meet = meets[ch-1]
        break
      } else {
        fmt.Print("Invalid Choice!\n\n")
      }
    }
  }
  return
}

func get_url(meet Meeting) (url string) {
  url = "zoommtg://zoom.us/join?confno=" + meet.Metno
  if meet.Paswd != "" {
    if runtime.GOOS == "windows" {
      url += "^"
    }
    url += "&pwd=" + meet.Paswd
  }
  return
}

func open(url string) error {
  var cmd string
  var args []string

  switch runtime.GOOS {
  case "windows":
    cmd = "cmd"
    args = []string{"/c", "start"}
  case "darwin":
    cmd = "open"
  default:
    cmd = "xdg-open"
  }
  args = append(args, url)
  return exec.Command(cmd, args...).Start()
}

func main() {
  config_file := "config.json"
  val, pres := os.LookupEnv("ZOOOM_CONFIG")
  if pres {
    config_file = val
  }
  config := read_config(config_file)
  meet := choose_meet(find_meets(config))
  if meet.Name == "" {
    fmt.Println("No Meeting found!")
    fmt.Print("Press Enter to Exit...")
    fmt.Scanln()
  } else {
    fmt.Println("Joining Meeting:", meet.Name)
    open(get_url(meet))
  }
}
