# Go-MultipleChoice [![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/thewolfnl/go-multiplechoice)
If you're building a CLI application in GO and require the user to make a selection out of given options, you can use this package.

![Demo](https://media.giphy.com/media/SiEr1FJcTEomgLal4o/giphy.gif)

## Install :package:
```
import (
    "github.com/thewolfnl/go-multiplechoice"
)
```

## Usage :radio_button:
```
// Single selection
selection := MultipleChoice.Selection("Select one: ", []string{"option1", "option2", "option3"}])

// Multi selection
selections := MultipleChoice.MultiSelection("Select one: ", []string{"option1", "option2", "option3"}])
```

## Suggestions :thought_balloon:
Please create an [issue](https://github.com/TheWolfNL/go-multiplechoice/issues) if you have a suggestion.
