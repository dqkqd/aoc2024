use std::{io::BufRead, str::FromStr};

use anyhow::Result;

use rust::read;

enum Op {
    Add,
    Mul,
    Concat,
}

#[derive(Debug)]
struct Equation {
    lhs: i64,
    rhs: Vec<i64>,
}

fn binary_ops_generator(size: usize) -> impl Iterator<Item = Vec<Op>> {
    (0..(1 << size)).map(move |mask| {
        (0..size)
            .map(|b| {
                if (mask & (1 << b)) > 0 {
                    Op::Add
                } else {
                    Op::Mul
                }
            })
            .collect()
    })
}

fn ternary_ops_generator(size: usize) -> impl Iterator<Item = Vec<Op>> {
    (0..(3_u32.pow(size as u32))).map(move |mask| {
        (0..size)
            .map(|b| {
                let pow = 3_u32.pow(b as u32);
                let s = (mask / pow) % 3;
                if s == 0 {
                    Op::Add
                } else if s == 1 {
                    Op::Mul
                } else {
                    Op::Concat
                }
            })
            .collect()
    })
}

impl FromStr for Equation {
    type Err = anyhow::Error;

    fn from_str(value: &str) -> Result<Self, Self::Err> {
        let lhs = value.split(":").next().unwrap_or_default().parse::<i64>()?;
        let rhs = value
            .split(":")
            .last()
            .unwrap_or_default()
            .split(" ")
            .map(|s| s.trim())
            .filter_map(|s| s.parse::<i64>().ok())
            .collect();
        Ok(Equation { lhs, rhs })
    }
}

impl Equation {
    fn calculate(&self, ops: Vec<Op>) -> i64 {
        let mut acc = self.rhs[0];
        for (value, op) in self.rhs.iter().skip(1).zip(ops.iter()) {
            match op {
                Op::Add => acc += value,
                Op::Mul => acc *= value,
                Op::Concat => {
                    let log = (*value as f64).log10().floor() as u32;
                    let pow = 10_u32.pow(log + 1);
                    acc = acc * (pow as i64) + value;
                }
            }
        }
        acc
    }
}

fn read_equations() -> Result<Vec<Equation>> {
    let buf = read(7, false)?;
    Ok(buf
        .lines()
        .map_while(Result::ok)
        .filter_map(|line| Equation::from_str(&line).ok())
        .collect())
}

fn part1() -> Result<()> {
    let sum = read_equations()?
        .iter()
        .filter(|equation| {
            binary_ops_generator(equation.rhs.len() - 1)
                .any(|ops| equation.calculate(ops) == equation.lhs)
        })
        .map(|equation| equation.lhs)
        .sum::<i64>();

    println!("{}", sum);
    Ok(())
}

fn part2() -> Result<()> {
    let sum = read_equations()?
        .iter()
        .filter(|equation| {
            ternary_ops_generator(equation.rhs.len() - 1)
                .any(|ops| equation.calculate(ops) == equation.lhs)
        })
        .map(|equation| equation.lhs)
        .sum::<i64>();

    println!("{}", sum);
    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;

    Ok(())
}
