# chooser

chooser is a library for creating command-line option pickers.

### Usage

```go
import "github.com/nickbarth/chooser"

func main() {
	fmt.Println("Pick a number:")

	ch := chooser.NewChooser(5, 20)
	choice := ch.Choose([]string{"one", "two", "three", "four", "five", "six"})

	fmt.Println("You Chose:", choice)
}
```

### License
WTFPL &copy; 2018 Nick Barth
