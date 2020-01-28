# Go-Enum

Go-Enum is a tool to automatigically generate String() and Parse() functions in generated files located next to the enum file.

# Commands
goenum --path \<PATH\> (--debug)

If debug mode is enabled, it will print which enums files have received treatment and their path.

# Example
For a file looking like the below code:
```
package relationstate

// x-enum
type RelationState int

const (
    Pending RelationState = iota
    Connected
)
```

The below code will me automagically be generated:
```
// Generated code. DO NOT EDIT.
package relationstate

import "errors"

func (kind RelationState) String() string {
    switch kind { 
    case 0:
        return "Pending" 
    case 1:
        return "Connected" 
    default:
        return ""
    }
}

func Parse(name string) (RelationState, error) {
    switch name { 
    case "Pending":
        return RelationState(0), nil 
    case "Connected":
        return RelationState(1), nil 
    default:
        return RelationState(0), errors.New("Enum for \"RelationState\" not found using name = " + name)
    }
}
```