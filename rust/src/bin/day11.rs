use std::{collections::HashMap, io::BufRead};

use anyhow::Result;

use rust::read;

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Hash)]
enum Stone {
    Zero,
    Odd { n: u64 },
    Even { n: u64, digits: u32 },
}

impl From<u64> for Stone {
    fn from(n: u64) -> Self {
        if n == 0 {
            Stone::Zero
        } else {
            let digits = ((n as f64).log10() as u32) + 1;
            if digits % 2 == 1 {
                Stone::Odd { n }
            } else {
                Stone::Even { n, digits }
            }
        }
    }
}

impl Stone {
    fn divide(&self) -> Vec<Stone> {
        match self {
            Stone::Zero => vec![Stone::Odd { n: 1 }],
            Stone::Odd { n } => {
                vec![Stone::from(n * 2024)]
            }
            Stone::Even { n, digits } => {
                let pow = 10u64.pow(digits / 2);
                vec![Stone::from(n / pow), Stone::from(n % pow)]
            }
        }
    }
}

fn read_stones() -> Result<Vec<Stone>> {
    let reader = read(11, false)?;
    let line = reader.lines().next().unwrap()?;
    let stones = line
        .split(" ")
        .filter_map(|s| s.parse::<u64>().ok())
        .map(Stone::from)
        .collect();
    Ok(stones)
}

fn run(times: u8) -> Result<()> {
    let stones = read_stones()?;
    let mut counts = HashMap::new();
    for stone in stones {
        counts.entry(stone).or_insert(1);
    }

    for _ in 0..times {
        let vec_counts: Vec<(Stone, u64)> = counts
            .drain()
            .flat_map(|(stone, count)| {
                stone
                    .divide()
                    .into_iter()
                    .map(|s| (s, count))
                    .collect::<Vec<(Stone, u64)>>()
            })
            .collect();
        for (stone, count) in vec_counts {
            *counts.entry(stone).or_insert(0) += count;
        }
    }
    let ans = counts.values().sum::<u64>();
    println!("{}", ans);
    Ok(())
}

fn part1() -> Result<()> {
    run(25)
}

fn part2() -> Result<()> {
    run(75)
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
