package main

import (
  "flag"
  "fmt"
  "os"
  "os/exec"
  "os/signal"
  "math/rand"
  "regexp"
  "strconv"
  "sync"
  "syscall"
  "time"

  "github.com/gomodule/redigo/redis"
)

var ADJECTIVES = []string{"afraid", "all", "angry", "beige", "big", "bitter", "blue", "brave", "breezy", "bright", "brown", "bumpy", "busy", "calm", "chatty", "chilled", "chilly", "chubby", "clean", "clever", "cold", "cool", "crazy", "cruel", "cuddly", "curly", "curvy", "cute", "cyan", "deep", "dirty", "dry", "dull", "eager", "early", "easy", "eight", "eighty", "eleven", "empty", "every", "evil", "fair", "famous", "fast", "few", "fifty", "five", "flat", "fluffy", "forty", "four", "fresh", "friendly", "fruity", "funny", "fuzzy", "gentle", "giant", "gold", "good", "gorgeous", "great", "green", "grumpy", "happy", "healthy", "heavy", "hip", "honest", "hot", "hungry", "icy", "itchy", "khaki", "kind", "large", "late", "lazy", "lemon", "light", "little", "long", "loose", "loud", "lovely", "lucky", "many", "mean", "metal", "mighty", "modern", "moody", "nasty", "neat", "nervous", "new", "nice", "nine", "ninety", "odd", "old", "olive", "orange", "perfect", "pink", "plenty", "polite", "poor", "popular", "pretty", "proud", "puny", "purple", "quick", "quiet", "rare", "real", "red", "rich", "rotten", "rude", "sad", "salty", "selfish", "serious", "seven", "shaggy", "shaky", "sharp", "shiny", "short", "shy", "silent", "silly", "silver", "six", "sixty", "slick", "slimy", "slow", "small", "smart", "smooth", "soft", "solid", "some", "sour", "spicy", "spotty", "stale", "strange", "strong", "stupid", "sweet", "swift", "tall", "tame", "tangy", "tasty", "ten", "tender", "thick", "thin", "thirty", "three", "tidy", "tiny", "tired", "tough", "tricky", "twelve", "twenty", "two", "unlucky", "upset", "vast", "violet", "warm", "weak", "wet", "wicked", "wide", "wild", "wise", "witty", "yellow", "young", "yummy"}
var NOUNS = []string{"apes", "animals", "bars", "baths", "breads", "bushes", "cloths", "clubs", "hoops", "loops", "memes", "papers", "parks", "paths", "showers", "signs", "sites", "streets", "teeth", "tires", "webs", "actors", "ads", "adults", "aliens", "ants", "apples", "apricots", "avocados", "baboons", "badgers", "bags", "balloons", "bananas", "bats", "beans", "bears", "beds", "beers", "bees", "berries", "bikes", "birds", "boats", "bobcats", "books", "bottles", "boxes", "brooms", "buckets", "bugs", "bulldogs", "buses", "buttons", "camels", "cameras", "candies", "candles", "carpets", "carrots", "carrots", "cars", "cats", "chairs", "cheetahs", "chefs", "cherries", "chicken", "clocks", "clouds", "coats", "cobras", "coins", "colts", "comics", "cooks", "cougars", "regions", "cows", "crabs", "crabs", "crews", "cups", "cycles", "dancers", "days", "deers", "dingos", "dodos", "dogs", "dolls", "dolphins", "donkeys", "donuts", "doodles", "doors", "dots", "dragons", "drinks", "dryers", "ducks", "ducks", "eagles", "ears", "eels", "eggs", "mammals", "emus", "experts", "eyes", "falcons", "fans", "feet", "files", "fireants", "fishes", "flies", "flowers", "forks", "foxes", "friends", "frogs", "games", "garlics", "geckos", "geese", "ghosts", "ghosts", "gifts", "glasses", "goats", "gorillas", "grapes", "guests", "hairs", "hats", "hornets", "horses", "hotels", "hounds", "houses", "humans", "icons", "impalas", "insects", "islands", "items", "jars", "jeans", "jobs", "jokes", "angaroos", "keys", "kids", "kings", "kiwis", "knives", "ladybugs", "lamps", "laws", "lemons", "lies", "ligers", "lions", "lizards", "llamas", "lobsters", "mails", "mangos", "maps", "masks", "mayflies", "meals", "melons", "mice", "mirrors", "moles", "monkeys", "months", "moons", "moose", "mugs", "nails", "needles", "news", "numbers", "olives", "onions", "oranges", "otters", "owls", "pandas", "pans", "panthers", "pants", "papayas", "parents", "parrots", "paws", "peaches", "pears", "peas", "penguins", "pens", "pets", "phones", "pianos", "pigs", "pillows", "planes", "planets", "plants", "plums", "poems", "poets", "points", "pots", "pugs", "pumas", "pumpkins", "queens", "rabbits", "radios", "rats", "ravens", "readers", "rice", "rings", "rivers", "rockets", "rocks", "roses", "rules", "schools", "scissors", "bats", "seals", "seas", "sheep", "shirts", "shoes", "shrimps", "singers", "sloths", "snails", "snakes", "socks", "spiders", "spies", "spoons", "squids", "stars", "steaks", "wings", "students", "suits", "suns", "swans", "symbols", "tables", "taxis", "teachers", "terms", "ties", "tigers", "timers", "tips", "toes", "tomatoes", "tools", "toys", "trainers", "trains", "trams", "trees", "turkeys", "turtles", "vans", "walls", "walls", "wasps", "waves", "ways", "weeks", "windows", "wolves", "wombats", "worms", "yaks", "years", "zebras", "zoos"}
var VERBS = []string{"accept", "act", "add", "admire", "agree", "allow", "appear", "applaud", "approve", "argue", "arrive", "attack", "attend", "bake", "bathe", "battle", "beam", "beg", "begin", "behave", "boil", "bow", "brake", "breathe", "brush", "build", "burn", "buy", "call", "camp", "care", "carry", "change", "cheat", "check", "cheer", "chew", "clap", "clean", "collect", "compare", "compete", "complain", "confess", "cough", "count", "cover", "crash", "cross", "cry", "dance", "decide", "deliver", "deny", "design", "destroy", "develop", "divide", "do", "double", "doubt", "draw", "dream", "dress", "drive", "drop", "drum", "eat", "end", "enjoy", "exercise", "exist", "explain", "explode", "fail", "fetch", "film", "fix", "flash", "float", "flow", "fly", "fold", "fry", "give", "glow", "grab", "greet", "grin", "grow", "guess", "hammer", "hang", "happen", "heal", "hear", "help", "hide", "hope", "hug", "hunt", "impress", "invent", "invite", "itch", "jam", "jog", "join", "joke", "judge", "juggle", "jump", "kick", "kiss", "kneel", "knock", "know", "laugh", "lay", "learn", "leave", "lick", "lie", "listen", "live", "look", "love", "march", "marry", "mate", "matter", "melt", "mix", "move", "nail", "notice", "obey", "occur", "own", "pay", "peel", "perform", "play", "poke", "press", "pretend", "promise", "protect", "prove", "provide", "pull", "pump", "punch", "push", "raise", "reflect", "refuse", "relate", "relax", "remain", "remember", "repair", "repeat", "reply", "report", "rescue", "rest", "retire", "return", "rhyme", "ring", "roll", "rule", "run", "rush", "scream", "search", "sell", "serve", "shake", "share", "shave", "shop", "shout", "sin", "sing", "sip", "sit", "sleep", "smash", "smell", "smile", "smoke", "sneeze", "sniff", "sort", "sparkle", "speak", "stare", "study", "suffer", "swim", "switch", "talk", "tan", "tap", "taste", "teach", "tease", "tell", "thank", "think", "tickle", "tie", "trade", "train", "travel", "try", "turn", "type", "unite", "vanish", "visit", "wait", "walk", "warn", "wash", "watch", "wave", "whisper", "wink", "wonder", "work", "worry", "yawn", "yell"}

type RedisCommand struct {
    commandName string
    args []interface{}
}

func main() {
  rand.Seed(time.Now().UnixNano())

  port := flag.Int("p", 8080, "Port to forward")
  subdomain := flag.String("s", "", "Subdomain to forward to. Leave blank for random.")
  isServer := flag.Bool("server", false, "Is running as server.")

  flag.Parse()

  if *isServer {
    os.Exit(Server(*port, subdomain))
  } else {
    os.Exit(Client(*port, subdomain))
  }
}

func SelectRandom(list []string) string {
  return list[rand.Intn(len(list))]
}

func RandomId() string {
  return fmt.Sprintf("%s-%s-%s", SelectRandom(ADJECTIVES), SelectRandom(NOUNS), SelectRandom(VERBS))
}

func c(commandName string, args ...interface{}) RedisCommand {
  return RedisCommand{
    commandName,
    args,
  }
}

func CleanupServer(client redis.Conn, router string, service string) {
  fmt.Println("Cleaning up...")

  err := DoMany(client,
    c("DEL", router + "rule"),
    c("DEL", router + "entrypoints/0"),
    c("DEL", router + "service"),
    c("DEL", router + "tls/certresolver"),
    c("DEL", router + "tls/domains/0/main"),
    c("DEL", router + "tls/domains/0/sans/0"),
    c("DEL", service + "loadBalancer/servers/0/url"),
  )

  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to clean up.")
  }
}

func DoMany(client redis.Conn, commands ...RedisCommand) error {
  for _, cmd := range commands {
    err := client.Send(cmd.commandName, cmd.args...)
    if err != nil {
      return err
    }
  }

  err := client.Flush()
  return err
}

func Server(port int, subdomain *string) int {
  if len(*subdomain) == 0 {
    fmt.Fprintln(os.Stderr, "Invalid subdomain.")
    return 1
  } else {
    match, err := regexp.MatchString("^[a-z0-9-]+$", *subdomain)
    if !match || err != nil {
      fmt.Fprintln(os.Stderr, "Invalid subdomain. Please include only lowercase alphanumeric characters, or a dash.")
      return 1
    }
  }

  if port < 10000 || port > 65535 {
    fmt.Fprintln(os.Stderr, "Invalid port.")
    return 1
  }

  router := fmt.Sprintf("traefik/http/routers/%s/", *subdomain)
  service := fmt.Sprintf("traefik/http/services/tunnel_%s/", *subdomain)

  client, err := redis.Dial("unix", "/home/pantherman594/traefik/traefik_redis/redis.sock")
  if err != nil {
    fmt.Fprintln(os.Stderr, "Failed to connect to redis.")
    return 1
  }
  defer client.Close()

  _, err = redis.String(client.Do("GET", router + "rule"))
  if err != redis.ErrNil {
    fmt.Println(fmt.Errorf("Subdomain is already in use."))
    return 1
  }

  defer CleanupServer(client, router, service)
  _, err = client.Do("SET", service + "loadBalancer/servers/0/url",fmt.Sprintf("http://host.docker.internal:%d", port))
  if err != nil {
    fmt.Fprintln(os.Stderr, "Failed to create treafik configuration.")
    return 1
  }
  err = DoMany(client,
    c("SET", router + "rule", fmt.Sprintf("Host(`%s.tt.dav.sh`)", *subdomain)),
    c("SET", router + "entrypoints/0", "websecure"),
    c("SET", router + "service", fmt.Sprintf("tunnel_%s", *subdomain)),
    c("SET", router + "tls/certresolver", "dnsresolver"),
    c("SET", router + "tls/domains/0/main", "tt.dav.sh"),
    c("SET", router + "tls/domains/0/sans/0", "*.tt.dav.sh"),
  )
  if err != nil {
    fmt.Fprintln(os.Stderr, "Failed to create treafik configuration.")
    return 1
  }

  fmt.Printf("Tunnel accessible at https://%s.tt.dav.sh\n", *subdomain)

  var wg sync.WaitGroup

  // Catch end
  signalChannel := make(chan os.Signal, 2)
  keyboardChannel := make(chan bool)
  signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
  wg.Add(1)
  go func() {
    defer wg.Done()

    for {
      select {
      case <-signalChannel:
        return
      case <-keyboardChannel:
        return
      }
    }
  }()

  // Read keyboard input for q
  go func() {
    // Disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // Do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    var b []byte = make([]byte, 1)

    fmt.Println("\nPress q to quit.")

    for {
      os.Stdin.Read(b)
      if b[0] == 'q' {
        keyboardChannel <- true
        return
      }
    }
  }()

  wg.Wait()
  return 0
}

func Client(port int, subdomain *string) int {
  // Pick a random port [10000, 65535]
  remotePort := 10000 + rand.Intn(55536)

  sub := *subdomain
  if len(sub) == 0 {
    sub = RandomId()
  } else {
    match, err := regexp.MatchString("^[a-z0-9-]+$", sub)
    if !match || err != nil {
      fmt.Fprintln(os.Stderr, "Invalid subdomain. Please include only lowercase alphanumeric characters, or a dash.")
      return 1
    }
  }

  fmt.Printf("Establishing tunnel to port %d...\n", port)

  cmd := exec.Command("ssh",
    "-tR", fmt.Sprintf(":%d:localhost:%d", remotePort, port), "culatra",
    "~/tunnel", "-server", "-p", strconv.Itoa(remotePort), "-s", sub)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Stdin = os.Stdin

  err := cmd.Run()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Command errored with", err)
  }

  return 0
}
