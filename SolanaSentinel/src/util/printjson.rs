use serde::Serialize;
use serde_json;

pub fn print_json<T: Serialize>(item: &T) {
    let json_str = serde_json::to_string_pretty(item).unwrap();
    println!("{}", json_str);
}