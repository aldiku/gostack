package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

func Log() *zap.Logger {
	if log != nil {
		return log
	}

	// File logger configuration
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // Set the log file name.
		MaxSize:    1,              // Maximum size in megabytes before rotation.
		MaxBackups: 2,              // Maximum number of old log files to retain.
		MaxAge:     2,              // Maximum number of days to retain old log files.
		Compress:   false,          // Whether to compress old log files.
	}

	// Console logger configuration
	consoleEncoderConfig := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "timestamp",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // Enable color for log levels
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	// Create the WriteSyncer for the log file
	fileWriteSyncer := zapcore.AddSync(lumberjackLogger)

	// Create the Zap core for the log file
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp" // Change the key name for timestamps.
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		fileWriteSyncer,
		zapcore.InfoLevel,
	)

	// Create the WriteSyncer for the console
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	// Create the Zap core for the console
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		consoleWriteSyncer,
		zapcore.InfoLevel,
	)

	// Use zap.New to create a logger with two cores, one for the file and another for the console
	log = zap.New(zapcore.NewTee(fileCore, consoleCore), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return log
}

type LogMessage struct {
	Level      string `json:"level"`
	Timestamp  string `json:"timestamp"`
	Caller     string `json:"caller"`
	Msg        string `json:"msg"`
	Stacktrace string `json:"stacktrace"`
}

func GetLogger(c echo.Context) error {
	key := c.QueryParam("key")
	mode := c.QueryParam("mode")
	if key != "nganu" {
		return c.JSON(401, "maaf anda tidak berhak")
	}
	if mode == "ajax" {
		limit := c.QueryParam("limit")
		ln := 50
		if limit != "" {
			ln, _ = strconv.Atoi(limit)
		}
		var logs []LogMessage
		latest, _ := ReverseReadLines("./logs/app.log", ln)
		for _, v := range latest {
			logEntry, err := parseLogLine(v)
			if err == nil {
				logs = append(logs, logEntry)
			}
		}
		return c.JSON(200, logs)
	} else {
		tmpl, err := template.New("example").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <style>
      body {
        font-family: "Courier New", monospace;
        background-color: #1a1a1a;
        color: #00ff00;
        margin: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
      }

      .terminal {
        background-color: #000;
        border: 2px solid #00ff00;
        padding: 20px;
        width: 80%;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        height: 90%;
        overflow: scroll;
      }

      .command {
        color: #00ff00;
      }

      .output {
        color: #00ff00;
      }
    </style>
    <title>Terminal-Like Box</title>
  </head>
  <body>
    <div class="terminal">
      <div class="output">Initializing...</div>
    </div>

    <script>
      function scrollToBottom() {
        const terminal = document.querySelector(".terminal");
        terminal.scrollTop = terminal.scrollHeight;
      }

      function updateTerminal() {
        fetch("{{.API_LOG}}")
          .then((response) => response.json())
          .then((data) => {
            const terminalOutput = document.querySelector(".terminal .output");
            terminalOutput.innerHTML = "";
            data.forEach((entry) => {
              const logLine = document.createElement("div");
			  var tim = entry.timestamp;
			  var lvl = entry.level;
			  var msgg = entry.msg;
              logLine.textContent = tim +" ["+ lvl +"] "+ msgg ;
              terminalOutput.appendChild(logLine);
            });
            scrollToBottom();
          })
          .catch((error) => {
            console.error("Error:", error);
            const terminalOutput = document.querySelector(".terminal .output");
            terminalOutput.textContent = "Error fetching data";
            scrollToBottom();
          });
      }

      setInterval(updateTerminal, 10000);

      updateTerminal();
    </script>
  </body>
</html>
`)
		if err != nil {
			fmt.Println(err)
		}
		// err := c.Render(200, htmlString, map[string]interface{}{
		// 	"API_LOG": os.Getenv("BASE_URL") + "/debug/console?key=nganu&mode=ajax",
		// })
		err = tmpl.Execute(c.Response().Writer, map[string]interface{}{
			"API_LOG": os.Getenv("BASE_URL") + "/debug/console?key=nganu&mode=ajax",
		})
		fmt.Println(err)
		return err
	}
}

func getTotalLines(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineNum, nil
}

func parseLogLine(line string) (logMessage LogMessage, err error) {
	err = json.Unmarshal([]byte(line), &logMessage)
	return
}

func ReverseReadLines(filename string, maxLines int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the total number of lines in the log file
	totalLines, _ := getTotalLines(file)

	// Calculate the starting line number
	startLine := totalLines - maxLines
	if startLine < 0 {
		startLine = 0
	}

	// Reset the read position of the file
	file.Seek(0, io.SeekStart)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	lineNum := 0
	var lines []string
	for scanner.Scan() {
		lineNum++

		if lineNum <= startLine {
			continue
		}

		line := scanner.Text()
		lines = append(lines, line)

		if len(lines) == maxLines {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
