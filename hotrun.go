// description       : An utilities tool build with GoLang to help user to run build certains files without enter a set of commands with flags (similar with Makefile)
//                     supported GoLang, Docker, C and etc...
// version           : "1.0.0"
// creator           : Gordon Lim <honwei189@gmail.com>
// created           : 25/09/2019 18:41:21
// last modified     : 06/01/2020 19:40:42
// last modified by  : Gordon Lim <honwei189@gmail.com>

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	// . "github.com/logrusorgru/aurora"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/hpcloud/tail"
	"github.com/urfave/cli"
)

//
var args = []string{}
var cliArgs = []string{}
var watcher *fsnotify.Watcher
var file string
var dir string

// var pause int = 0

func runCommand(filename string) {

	// args = append(args, filename)

	// for index, each := range os.Args {
	// 	if index > 1 {
	// 		args = append(args, each)
	// 	}
	// }

	shellCmd := ""
	extension := ""
	clsCmd := ""

	if runtime.GOOS == "windows" {
		clsCmd = "cls"
	} else {
		clsCmd = "clear"
	}

	switch strings.ToLower(filename) {
	case "makefile":
		extension = "cpp"
	case "dockerfile":
		extension = "docker"
	default:
		extension = filepath.Ext(filename)
	}

	extension = strings.ToLower(strings.Replace(extension, ".", "", 1))

	switch extension {
	// case "php":
	// 	shellCmd = "php"
	// case "sh":
	// 	shellCmd = "sh"
	case "c", "cpp":
		cmdRun2("gcc", "", filename)
	case "docker":
		// for _, each := range cliArgs {
		// 	filename = filename + " " + each
		// }

		args = append(args, "build")
		args = append(args, "-t")

		for _, each := range cliArgs {
			args = append(args, each)
		}

		// fmt.Println(args)
		cmdRun3("docker", args)
		args = nil
	case "go":
		cmd := exec.Command(clsCmd) //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()

		cmdRun2("go", "run", filename)

	case "conf", "txt", "log", "sql", "md":
		if runtime.GOOS == "windows" {
			// t, _ := tail.TailFile(filename, tail.Config{Follow: true, Location: &tail.SeekInfo{0, 2}, Logger: tail.DiscardingLogger})
			// Location:  &tail.SeekInfo{Offset: 0, Whence: 4},
			// Location:  &tail.SeekInfo{0, 4},
			// cmdRun2("cmd", "/c", "cls")
			cmd := exec.Command(clsCmd) //Windows example, its tested
			cmd.Stdout = os.Stdout
			cmd.Run()

			t, _ := tail.TailFile(filename, tail.Config{
				ReOpen:    true,
				Follow:    true,
				MustExist: false,
				Poll:      true,
				Logger:    tail.DiscardingLogger,
				Location:  &tail.SeekInfo{Offset: 0, Whence: 1},
			})

			for line := range t.Lines {
				fmt.Println(line.Text)
			}
		} else {
			cmd := exec.Command(clsCmd) //Windows example, its tested
			cmd.Stdout = os.Stdout
			cmd.Run()

			cmdRun("tail", args, filename)
		}

	default:
		cmd := exec.Command(clsCmd) //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()

		shellCmd = extension
		cmdRun(shellCmd, args, filename)
		break
	}

}

func cmdRun(shellCmd string, args []string, filename string) {
	args = append(args, filename)

	for _, each := range cliArgs {
		args = append(args, each)
	}

	// log.Println(args)

	if len(shellCmd) > 0 {
		cmd := exec.Command(shellCmd, args...)
		args = nil

		// create a pipe for the output of the script
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			os.Exit(0)
			return
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				// fmt.Printf("\t > %s\n", scanner.Text())
				// println(scanner.Text())
				fmt.Printf("%s\n", scanner.Text())
			}
		}()

		// bufio.NewReaderSize(cmdReader, 20000000000)
		// scanner := bufio.NewScanner(cmdReader)
		// go func() {
		// 	// buf := make([]byte, 0, 64*1024)
		// 	// scanner.Buffer(buf, 10240*1024*1024)

		// 	const maxCapacity = 512 * 8096
		// 	buf := make([]byte, maxCapacity)
		// 	scanner.Buffer(buf, maxCapacity*(8192*8192)*256)
		// 	for scanner.Scan() {
		// 		// fmt.Printf("\t > %s\n", scanner.Text())
		// 		// println(scanner.Text())
		// 		fmt.Printf("%s\n", scanner.Text())
		// 	}
		// }()

		// scanner := bufio.NewScanner(cmdReader)
		// // scanner.Split(bufio.ScanWords)
		// count := 0
		// go func() {
		// 	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// 		count++
		// 		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
		// 		return 0, nil, nil
		// 	}
		// 	scanner.Split(split)
		// 	buf := make([]byte, 1024*8096)
		// 	scanner.Buffer(buf, bufio.MaxScanTokenSize*(8192*8192)*256)
		// 	for scanner.Scan() {
		// 		// fmt.Printf("%s\n", scanner.Text())
		// 	}

		// 	buf = nil
		// 	scanner = nil
		// 	println(count)
		// }()
		err = cmd.Start()
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			fmt.Fprintln(os.Stderr, "\n\n"+color.FgRed.Render(err)+"\n")
			os.Exit(0)
			return
		}

		err = cmd.Wait()
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
			fmt.Fprintln(os.Stderr, "\n\n"+color.FgRed.Render(err)+"\n")
			os.Exit(0)
			return
		}
	}
}

func cmdRun2(shellCmd string, opt string, filename string) {
	cmd := exec.Command(shellCmd, opt, filename)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()
	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()
	err := cmd.Wait()
	if err != nil {
		// log.Fatalf("cmd.Run() failed with %s\n", err)
		// fmt.Println("\n\nCommand error... Please confirm your confirm")
		// fmt.Println("\n\n" + color.FgRed.Render("Command error... Please confirm") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	if errStdout != nil || errStderr != nil {
		// log.Fatalf("failed to capture stdout or stderr\n")
		// fmt.Println("\n\n" + color.FgRed.Render("Unable to capture output from command...") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	// outStr, errStr := string(stdout), string(stderr)
	// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	// outStr, _ := string(stdout), string(stderr)
	// // fmt.Println(outStr)
	// outStr = ""
	stdoutIn = nil
	stderrIn = nil
	errStdout = nil
	stderr = nil
	stdout = nil
}

func cmdRun3(shellCmd string, opt []string) {
	cmd := exec.Command(shellCmd, opt...)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()
	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()
	err := cmd.Wait()
	if err != nil {
		// log.Fatalf("cmd.Run() failed with %s\n", err)
		// fmt.Println("\n\nCommand error... Please confirm your confirm")
		// fmt.Println("\n\n" + color.FgRed.Render("Command error... Please confirm") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	if errStdout != nil || errStderr != nil {
		// log.Fatalf("failed to capture stdout or stderr\n")
		// fmt.Println("\n\n" + color.FgRed.Render("Unable to capture output from command...") + "\n")
		fmt.Println("\n\n" + color.FgRed.Render("Terminated ...") + "\n")
		os.Exit(0)
	}
	// outStr, errStr := string(stdout), string(stderr)
	// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	// outStr, _ := string(stdout), string(stderr)
	// fmt.Println(outStr)
	// outStr = ""
	stdoutIn = nil
	stderrIn = nil
	errStdout = nil
	stderr = nil
	stdout = nil
}

// func readLines(reader string, out <-chan string) error {
// 	scanner := bufio.NewScanner(reader)
// 	const maxCapacity = 512 * 1024
// 	buf := make([]byte, maxCapacity)
// 	scanner.Buffer(buf, maxCapacity)

// 	for scanner.Scan() {
// 		out <- scanner.Text()
// 	}
// 	return scanner.Err()

// }

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			os.Stdout.Write(d)
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
	// never reached
	// panic(true)
	// return nil, nil
}

func monitoring(watcher *fsnotify.Watcher, ff string, typ string) {
	if typ == "d" {
		// paths := []string{
		// 	"D:\\AES Encrypt\\Deb",
		// 	"D:\\AES Encrypt\\M",
		// }

		paths := []string{}

		err := filepath.Walk(ff, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				paths = append(paths, path)
			}
			return nil
		})

		if err != nil {
			fmt.Println("ERROR", err)
		}

		for _, path := range paths {
			// log.Printf("Watching: %s\n", path)
			err := watcher.Add(path)
			if err != nil {
				log.Fatalf("Failed to watch directory: %s", err)
			}
		}

		paths = nil

		if len(file) > 0 {
			runCommand(file)
		}
	} else if typ == "f" {
		err := watcher.Add(ff)
		if err != nil {
			log.Fatal(err)
		} else {
			runCommand(ff)
		}
	}

	wg := &sync.WaitGroup{}

	go func() {
		for {
			watchUpdate(wg, watcher)
		}
	}()
	wg.Wait()
}

func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory

	if fi.IsDir() {
		watcher.Add(path)
	}

	return nil
}

func watchUpdate(wg *sync.WaitGroup, watcher *fsnotify.Watcher) {
	var extension = ""
	var filename = ""
	// var action = ""
	var curr int
	var last int
	var idx int

	for {
		select {
		case e := <-watcher.Events:
			extension = filepath.Ext(e.Name)
			extension = strings.Replace(extension, ".", "", 1)
			filename = e.Name
			// action = e.Op.String()
			// log.Println("修改文件：" + e.Name)
			// log.Println("修改類型：" + e.Op.String())
			idx++
		case err := <-watcher.Errors:
			log.Printf("Watcher error: %s\n", err.Error()) // No need to exit here
			// case <-time.After(time.Second * 2):
			// println("AA")
		}

		curr = curr + 1
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			time.AfterFunc(time.Second*1, func() {
				// log.Printf(":%d %d", idx, i)

				if runtime.GOOS == "windows" {
					if (i + 1) == idx {
						if len(file) > 0 && len(dir) > 0 {
							runCommand(file)
						} else {
							runCommand(filename)
						}

						idx = 0
						curr = 0
						last = 0
					}
				} else {
					if (i + 2) == idx {
						// log.Printf(":%s", filename)
						// cmd := exec.Command("php", filename, "--user=rosenaha --sdate=2018-10-20 --edate=2019-09-19")
						// if extension == "php" {
						if len(file) > 0 && len(dir) > 0 {
							runCommand(file)
						} else {
							runCommand(filename)
						}
						// }

						idx = 0
						curr = 0
						last = 0
					}
				}
			})

			last++
			wg.Done()
		}(wg, last)

		// if idx == last {
		// 	idx = 0
		// 	// curr = 0
		// 	log.Printf(":%s %d %d", action, curr, last)
		// }

		// time.After(time.Second * 5)
		// break
	}

	// println("Break")
}

func main() {
	var ff string
	var typ string
	// var numFlags int

	app := cli.NewApp()
	// app.Name = "Monitoring PHP CLI run changes"
	app.Name = "hotrun"
	app.EnableBashCompletion = true
	app.Usage = "\n\n\t\t\tHot run / reload / build file without manual run command likes \"php MY_FILE.php\""
	app.UsageText = app.Name + " [global options] command [command options] [arguments...]\n\n\t\t\texample:\n\n\t\t\t" + app.Name + " -f migrate.php run --user=my_user --sdate=2018-10-20 --edate=2019-09-19\n\n\t\t\t" + app.Name + " -f test.sh\n\n\t\t\t" + app.Name + " -f test.sh arg_1 arg_2 arg_3\n\n\t\t\t" + app.Name + " -f test.go\n\n\t\t\t" + app.Name + " -f test.python\n\n\t\t\t" + app.Name + " -d /home/abc"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Gordon Lim",
			Email: "gordon@weki.com.my",
		},
	}
	app.Copyright = "2019 Gordon Lim\n"
	// app.ArgsUsage = "[args and such]"

	// app.Commands = []cli.Command{
	// 	cli.Command{
	// 		Name:        "doo",
	// 		Aliases:     []string{"do"},
	// 		Category:    "motion",
	// 		Usage:       "do the doo",
	// 		// UsageText:   "doo - does the dooing",
	// 		Description: "no really, there is a lot of dooing to be done",
	// 	},
	// }

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, f",
			Value:       "",
			Usage:       "File name",
			Destination: &file,
		},
		cli.StringFlag{
			Name:        "dir, d",
			Value:       "",
			Usage:       "Directory path",
			Destination: &dir,
		},
		// cli.IntFlag{
		// 	Name:        "wait, w",
		// 	Value:       0,
		// 	Usage:       "Seconds.  Wait for how many seconds and then hot re-run",
		// 	Destination: &pause,
		// },
	}

	app.Action = func(c *cli.Context) error {
		// numFlags = c.NumFlags()
		cliArgs = c.Args()

		if c.NumFlags() == 0 {
			cli.ShowAppHelp(c)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	if len(file) > 0 {
		ff = file
		typ = "f"
	} else if len(dir) > 0 {
		ff = dir
		typ = "d"
	}

	if len(file) > 0 && len(dir) > 0 {
		ff = dir
		typ = "d"
	}

	if ff != "" {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("Failed to create watcher: %s", err)
			// fmt.Println(Bold(Cyan("Failed to create watcher")))
		}
		defer watcher.Close()

		exit := make(chan bool)

		monitoring(watcher, ff, typ)

		<-exit // 用來 阻塞應用不退出，只能通過「殺死進程」的方式退出，如 按住 Ctrl + C 快捷鍵強制推出
		runtime.Goexit()
	}
}
