# scheduler

Scheduler for periodic/repetitive tasks execution under allowed ranges of time.

Created by [Marcos Ortega](https://mortegam.com/) for automatic database backups, notifications and monitoring. You can use it or modify it to your needs.

# Features

- Coded in Go.
- Schedule periodic taks.
- Define how often each task will be eavluated.
- Define min/max age of last succesfull execution of each task.
- Define range of allowed execution of each task.
- Define commands:
    - optionally each command can have a 'catch' command to be executed if the main command fails.
    - optionally each command can have a 'deferred' command to be executed if the main command succeeds.

# Compile it (Windows, Linux, Mac)

```
cd scheduler
go build cmd/scheduler/scheduler.go
```

# Run it (Windows, Linux, Mac)

Verify your configuration file:

```
./scheduler -c myConfig.json
```

Execute scheduler:

```
./scheduler -c myConfig.json -r
```

# Configuration

myConfig.json

```
{
	"client": {
		"configs": {
			"path": "/etc/scheduler/tasks"
			, "secsBetweenSync": 60
		}
		, "state": {
			"path": "/etc/scheduler/state.json"
			, "secsBetweenSave": 60
		}
	}
}
```

This configuration file:

- defines '/etc/scheduler/tasks' as the folder that contains the tasks (.json) to be executed.
- each 60 seconds the '/etc/scheduler/tasks' folder will be read to detect added/updated/removed tasks.
- the persistent-state will be loaded/saved into '/etc/scheduler/state.json'
- the persistente-state will be saved 60 seconds after it changes.

# Task

myEveryMinuteTask.json

```
{
	"version": 0
	, "timing": {
        "range": {
		       "min": "5s"
		       , "max": "10s"
        }
        , "age": {
		       "tick": "5s"
		       , "min": "30s"
           , "max": "1h"
	    }
	}
	, "commands": [
        {
            "execute": "echo 'Hi this minute!'"
            , "catch": "echo 'Error!'"
            , "deferred": "echo 'Success!'"
        }
    ]
}
```

Which:

 - executes between the 5th and 10th second (inclusive) of each minute.
 - is evaluated every 5 seconds (tick).
 - new executions must be 30s or older from last succesfull execution.
 - execution is forced if last succesfull execution is 1h or older.

myEveryYearTask.json

{
	"version": 0
	, "timing": {
        "range": {
		    "min": "January 0h"
		    , "max": "January 3h"
        }
        , "age": {
		      "tick": "1d"
		      , "min": "3M"
          , "max": "6M"
	    }
	}
	, "commands": [ ...]
}

Which:

 - executes once-a-year between the midnight and 3rd hour (inclusive) at january.
 - is evaluated every day (tick).
 - new executions must be 3 month or older from last succesfull execution.
 - execution is forced if last succesfull execution is 6 months or older.

# Contact

Visit [mortegam.com](https://mortegam.com/) for other projects.

May you be surrounded by passionate and curious people. :-)



