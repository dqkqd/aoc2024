use std::{collections::BTreeMap, error::Error, io::BufRead};

use rust::read;

fn read_input() -> Result<(Vec<i32>, Vec<i32>), Box<dyn Error>> {
    let buf = read(1, false)?;
    let res = buf
        .lines()
        .map(|line| {
            let items: Vec<i32> = line
                .expect("line must be read")
                .split_whitespace()
                .map(|s| s.parse::<i32>().expect("parse success"))
                .take(2)
                .collect();
            (items[0], items[1])
        })
        .unzip();

    Ok(res)
}

fn part1() -> Result<(), Box<dyn Error>> {
    let (mut arr1, mut arr2) = read_input()?;
    arr1.sort();
    arr2.sort();
    let ans = arr1
        .into_iter()
        .zip(arr2)
        .map(|(v1, v2)| (v1 - v2).abs() as i64)
        .sum::<i64>();

    println!("{}", ans);
    Ok(())
}

fn part2() -> Result<(), Box<dyn Error>> {
    let (arr1, arr2) = read_input()?;
    let map = arr2.into_iter().fold(BTreeMap::new(), |mut map, e| {
        *map.entry(e).or_insert(0) += 1;
        map
    });
    let ans = arr1
        .into_iter()
        .map(|v| (v as i64) * map.get(&v).unwrap_or(&0))
        .sum::<i64>();

    println!("{}", ans);
    Ok(())
}

fn main() -> Result<(), Box<dyn Error>> {
    part1()?;
    part2()?;
    Ok(())
}
