package main

import (
    "../deps/lessgo/data/rdc"
    "../deps/lessgo/pass"
    "../deps/lessgo/utils"
    "./conf"
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "os/user"
    "regexp"
    "strings"
    "time"
)

var (
    err          error
    cfg          conf.Config
    flagPrefix   = flag.String("prefix", "", "the prefix folder path")
    flagUserSet  = flag.Bool("userset", false, "Create a System Administrator")
    flagUserDel  = flag.Bool("userdel", false, "Delete a System Administrator")
    emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")
)

const (
    //CMDC_DEFAULT   = "\033[m"
    //CMDC_BLACK     = "\033[30m"
    CMDC_RED   = "\033[31m"
    CMDC_GREEN = "\033[32m"
    CMDC_BROWN = "\033[33m"
    //CMDC_BLUE      = "\033[34m"
    //CMDC_MAGANTA   = "\033[35m"
    //CMDC_CYAN      = "\033[36m"
    //CMDC_LIGHTGRAY = "\033[37m"
    CMDC_CLOSE = "\033[0m"
)

func main() {

    //
    if u, err := user.Current(); err != nil || u.Uid != "0" {
        //log.Fatal("Permission Denied : must be run as root")
    }

    //
    flag.Parse()
    if cfg, err = conf.NewConfig(*flagPrefix); err != nil {
        log.Fatal(err)
    }

    if cn, err := cfg.DatabaseInstance(); err == nil {
        rdc.InstanceRegister("def", cn)
    } else {
        log.Fatal(err)
    }

    if *flagUserSet {
        cmdUserSet()
    } else if *flagUserDel {
        //cmdUserDel()
    } else {
        fmt.Println(CMDC_RED + "No Command Found" + CMDC_CLOSE)
    }
}

func cmdUserSet() {

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Panic:", r)
        }
        exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
    }()

    // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak").Run()

    fmt.Println(CMDC_GREEN + "This wizard will guide you to create a System Administrator." + CMDC_CLOSE)

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        fmt.Println("Internal Server Error: Can not connect to database")
        os.Exit(1)
    }

    var email string
    for {

        email = ""

        fmt.Printf(CMDC_BROWN + "\nEnter a Email to login: " + CMDC_CLOSE)
        fmt.Scanf("%s", &email)

        email = strings.ToLower(strings.TrimSpace(email))
        if matched := emailPattern.MatchString(email); !matched {
            fmt.Printf(CMDC_RED + "Email is not valid, Please choose another one" + CMDC_CLOSE)
            continue
        }

        q := rdc.NewQuerySet().From("ids_login").Limit(1)
        q.Where.And("email", email)

        rsu, err := dcn.Query(q)
        if err == nil && len(rsu) == 1 {
            fmt.Printf(CMDC_RED + "The Email already exists, please choose another one" + CMDC_CLOSE)
            continue
        }

        break
    }

    //
    passwd := ""
    for {

        prompt := CMDC_BROWN + "\rEnter new password: " + CMDC_CLOSE
        reader := bufio.NewReaderSize(os.Stdin, 1)
        fmt.Printf(prompt)
        for {

            c, _ := reader.ReadByte()
            if c == '\n' {
                break
            }

            passwd += string(c)

            prompt += "*"
            fmt.Printf(prompt)
        }

        if len(passwd) >= 12 && len(passwd) <= 50 {
            break
        }

        fmt.Println(CMDC_RED + "Password must be between 12 and 50 characters long. Please choose another one" + CMDC_CLOSE)
    }
    hash, _ := pass.HashDefault(passwd)

    uname := utils.StringNewRand36(8)
    item := map[string]interface{}{
        "uname":   uname,
        "email":   email,
        "pass":    hash,
        "name":    uname,
        "status":  1,
        "group":   "1",
        "created": time.Now().Format(time.RFC3339), // TODO
        "updated": time.Now().Format(time.RFC3339), // TODO
    }
    rs, err := dcn.Insert("ids_login", item)
    if err != nil {
        fmt.Println("Internal Server Error: Can not write to database 2", err)
        os.Exit(1)
    }

    lastid, _ := rs.LastInsertId()
    itemgu := map[string]interface{}{
        "gukey":   fmt.Sprintf("%v.%v", lastid, 1),
        "created": time.Now().Format(time.RFC3339),
    }
    rs, err = dcn.Insert("ids_group_users", itemgu)
    if err != nil {
        fmt.Println("Internal Server Error: Can not write to database 3")
        os.Exit(1)
    }

    //
    fmt.Println(CMDC_GREEN + "Successfully created" + CMDC_CLOSE)
}

/*
func cmdUserDel() {

    //
    if *flagUserDel == "" {
        log.Fatal("Username can not be null")
    }
    email := strings.ToLower(*flagUserDel)



    fmt.Println("User deleted successfully")
}
*/