# Gocui

Gocui is a simple command line graphics toolkit for Go.Use it to build simple command line applications easily.

// At present, this is just a simple small project, welcome to help improve it.

# Features
- Easy to use. Just create an object and set its style, then call `Run()` to start the application.
- Customizable style. You can choose from different tokens to customize the style of the objects.
- Compatible. It use CSI codes to control the terminal.

# Functions
- Progress bar: Create a progress bar or an uncertain progress bar. And you can set the style of the progress bar.
- Text box: Create a text box to contain text.
- graph: Draw lines or curves in the terminal.

# Examples

## Progress bar
```go
p, _ := progress_bar.NewProgressBar("[%bar] %current/%total-%percent %rate", func(p *progress_bar.Property) {
		p.Style.BarComplete = "@"
		p.Style.BarIncomplete = "-"
	})
p.Run(time.Millisecond * 30)
// wait
<-p.Done
```

This will create a progress bar and run it.
Main goroutine need to wait for the progress bar to finish.

## Uncertain progress bar

```go
up, _ := progress_bar.NewProgressBar("[%bar] testing ubar...", func(p *progress_bar.Property) {
		p.Uncertain = true
		p.Style.BarIncomplete = " "
		p.Style.UnCertain = "<->"
	})
up.Run(time.Millisecond * 100)
// wait 5s
time.Sleep(time.Second * 5)
up.Stop()
```

This will create an uncertain progress bar and run it, then wait 5s and stop it.

## Text box
```go
payload := []string{
		"          Books Management System",
		"",
		" 1.Store new books    2.New user registration",
		" 3.Borrow books       4.Return books",
		" 5.All books          6.All user",
		" 7.Delete database    8.Log out",
		"",
		"          Select operation number:",
	}
window.ClearScreen()
aBox, _ := box.GetBox(len(payload)+2, 50+2, "bold", payload)
box.SetBoxAt(aBox, 0, 0)
```

This will create a text box and set it to the top left corner of the screen.

# Customization

// Currently, only Progress bar is supported.

Use a format string to customize the style of the objects,in which you can use the tokens to customize the style of the objects.

You can think of tokens as verbs in go.

## Tokens

Tokens that users use to customize the style of the objects.

### progress bar
- `%bar`: the progress bar
- `%current`: the current value
- `%total`: the total value
- `%percent`: the percentage
- `%elapsed`: the elapsed time
- `%rate`: Interval between two updates

# TODO
- [ ] Add more examples
- [ ] Add more modules
- [ ] Support more tokens
- [ ] Allow users define their own tokens
- [ ] Expand application scenarios, such as support parameters processing...
