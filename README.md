# Shake Programming Language


## Scope
#### A scope will always return
```go
x = {} // will by of type: empty
x = { return 1 } // will by of type: int, with value: 1
x = 1 // will by of type: int, with value: 1
x = if true { return 1 } // will by of type: int, with value: 1
```

## Conditionals
```go
if x == 1 {
    true {};
    false {};
}

if x == 1 {
    return 0;
}

if {
    x == 1 {};
    y == 1 {};
}

if x {
    == 1 {};
    >= 2 {};
}

if x == {
    1 {};
    2 {};
}
```

## Variables
```go
x: int = 1;
x = 1;
x = if 1 {
    true: 10 ;
    false { return 9 };
}
```

## Functions
```go
hello(): int
hello(): int {
    return 1;
}

// inline function
add(x: int, y: int) x + y
add(x: int, y: int): int: x + y

add(x: int, y: int): int {
    return if CheckIsEven(x, y) == {
        true: x + y + 10
        false: x + y + 5
    }
}
```

## Imports
```go
import "std/math"
import "IsEven"
import "github.com/ShakedGold/IsOdd"

import (
    math "std/math"
    "IsEven"
    "IsEven.shk"
    "github.com/ShakedGold/IsOdd"
)
```

## Exports
```js
export (
    CheckIsEven
    CheckIsEvenV2 as oops
)
```

## Structures
```rust
struct Person {
    Age: int: if Age > 18
    Name: string: if Age > 32 && Name.len > 2 || Name.len > 10
    pub Job: string
}

(p: Person) Hello(): string {
    return StringFormat("Hello: %", p.Name)
}

(entry) MainProg() {
    p1, err = Person {
        Age = 12
        Name = "S" // will error because of the constraint
    }

    // init to empty
    p2: Person

    p3, _ = Person {
        Age = 12
        Name = "Shaked"
        Job = "Clown"
    }

    // p3.Age panics
    // doesn't panic
    p3.Job

    // error code (implicit return 0)
    return if err == {
        /*
        return -1 if there is an error (err != empty)
        else return 0
        */
        empty: 0
        else: -1
    }
}
```

## Example Program
```go
import (
    "std/math"
    io "std/io"
    "IsEven"
)

x = 10
y: int = 10

// init to 0
z: int

CheckIsEven(x: int, y: int): int {
    a: int = if IsEvenV1(x) == {
        true: x + y
        false: {
            return 2 * x + Pow(y, 2)
        }
    }
    b: int = if IsEvenV2(y) == {
        true: y
        false: 2 * y
    }

    return a + b
}

// Entry decorator (can use multiple decorators)
(entry)
Main() {
    io.Print("Result: ", CheckIsEven(x, y), "\n")
}
```
