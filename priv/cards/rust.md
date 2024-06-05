# Rust

## Function with generic type

```
fn largest<T>(list: &[T]) -> T
```

## Create new vector

```
let v: Vec<i32> = Vec::new();
let v = vec![1, 2, 3];
```

## Create a hash map

```
use std::collections::HashMap;

let mut scores = HashMap::new();
scores.insert(String::from("Blue"), 10);
scores.insert(String::from("Yellow"), 50);
```

## Derive Debug

```
#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

println!("rect1 is {:?}", rect1);
println!("rect1 is {:#?}", rect1);
```

## Enum

```
enum IpAddr {
    V4(String),
    V6(String),
}

let home = IpAddr::V4(String::from("127.0.0.1"));
let loopback = IpAddr::V6(String::from("::1"));
```
