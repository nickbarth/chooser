# chooser

chooser is a library for creating command-line option pickers.

### Usage

```go
import "github.com/nickbarth/chooser"

func main() {
	ch := NewChooser(5, 20)
	choice := chr.Choose([]string{"one", "two", "three", "four", "five", "six"})

	fmt.Println("You Chose:", choice)
}
```

### License
WTFPL &copy; 2018 Nick Barth
