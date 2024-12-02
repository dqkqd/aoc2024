use std::{error::Error, io::BufRead};

use rust::read;

fn safe(v: &[i32]) -> bool {
    let is_diff = |v: &[i32]| {
        let s = (v[0] - v[1]).abs();
        (1..=3).contains(&s)
    };
    let is_incr = |v: &[i32]| v[0] < v[1];
    let is_decr = |v: &[i32]| v[0] > v[1];
    v.windows(2).all(is_diff) && (v.windows(2).all(is_incr) || v.windows(2).all(is_decr))
}

fn safe_rm(v: &[i32]) -> bool {
    let skip = |v: &[i32], i: usize| -> Vec<i32> {
        v[..i].iter().chain(v[i + 1..].iter()).cloned().collect()
    };

    safe(v) || (0..v.len()).map(|i| skip(v, i)).any(|v| safe(&v))
}

fn run(f: fn(&[i32]) -> bool) -> Result<(), Box<dyn Error>> {
    let buf = read(2, false)?;

    let safe_cnt: i32 = buf
        .lines()
        .map(|line| {
            let levels: Vec<i32> = line
                .expect("line must be read")
                .split_whitespace()
                .map(|s| s.parse::<i32>().expect("parse success"))
                .collect();
            if f(&levels) {
                1
            } else {
                0
            }
        })
        .sum();

    println!("{}", safe_cnt);
    Ok(())
}

fn part1() -> Result<(), Box<dyn Error>> {
    run(safe)
}

fn part2() -> Result<(), Box<dyn Error>> {
    run(safe_rm)
}

fn main() -> Result<(), Box<dyn Error>> {
    part1()?;
    part2()?;
    Ok(())
}
