use std::{
    collections::BTreeSet,
    io::BufRead,
    ops::{Add, Mul, Sub},
};

use anyhow::Result;

use rust::read;

fn cartesian_product<T: Clone, U: Clone>(ts: &[T], us: &[U]) -> Vec<(T, U)> {
    ts.iter()
        .flat_map(|t| {
            us.iter()
                .map(|u| (t.clone(), u.clone()))
                .collect::<Vec<(T, U)>>()
        })
        .collect()
}

#[derive(Default, Debug)]
struct Map {
    data: Vec<Vec<char>>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
struct Position(i32, i32);

impl Add<&Position> for &Position {
    type Output = Position;

    fn add(self, rhs: &Position) -> Self::Output {
        Position(self.0 + rhs.0, self.1 + rhs.1)
    }
}

impl Sub<&Position> for &Position {
    type Output = Position;

    fn sub(self, rhs: &Position) -> Self::Output {
        Position(self.0 - rhs.0, self.1 - rhs.1)
    }
}

impl Mul<i32> for &Position {
    type Output = Position;

    fn mul(self, rhs: i32) -> Self::Output {
        Position(self.0 * rhs, self.1 * rhs)
    }
}

impl Map {
    fn read() -> Result<Map> {
        let reader = read(8, false)?;
        let mut data = vec![];
        for line in reader.lines() {
            let line = line?;
            let line = line.trim();
            if !line.is_empty() {
                data.push(line.chars().collect())
            }
        }
        Ok(Map { data })
    }

    fn unique_atenna(&self) -> Vec<char> {
        let mut unique_chars = BTreeSet::new();
        for cs in &self.data {
            for c in cs {
                if c != &'.' {
                    unique_chars.insert(*c);
                }
            }
        }
        unique_chars.into_iter().collect()
    }

    fn _pairwise_positions(&self, antenna: &char) -> Vec<(Position, Position)> {
        let locations: Vec<Position> = self
            .data
            .iter()
            .enumerate()
            .flat_map(|(x, cs)| {
                cs.iter()
                    .enumerate()
                    .filter(|&(_, c)| c == antenna)
                    .map(|(y, _)| Position(x as i32, y as i32))
                    .collect::<Vec<Position>>()
            })
            .collect();
        cartesian_product(&locations, &locations)
    }

    fn pairwise_positions(&self) -> Vec<(Position, Position)> {
        self.unique_atenna()
            .iter()
            .flat_map(|c| self._pairwise_positions(c))
            .collect()
    }

    fn height(&self) -> usize {
        self.data.len()
    }

    fn width(&self) -> usize {
        self.data[0].len()
    }

    fn all_positions(&self) -> Vec<Position> {
        let xs = (0..self.height()).collect::<Vec<usize>>();
        let ys = (0..self.width()).collect::<Vec<usize>>();
        cartesian_product(&xs, &ys)
            .into_iter()
            .map(|p| Position(p.0 as i32, p.1 as i32))
            .collect()
    }
}

fn run(good_position: fn(&(Position, Position), &Position) -> bool) -> Result<()> {
    let map = Map::read()?;
    let pairwise_positions = map.pairwise_positions();
    let cnt = map
        .all_positions()
        .iter()
        .filter(|p| pairwise_positions.iter().any(|ps| good_position(ps, p)))
        .count();
    println!("{}", cnt);
    Ok(())
}

fn part1() -> Result<()> {
    run(|(p1, p2): &(Position, Position), p: &Position| -> bool {
        p1 != p2 && ((p + p1) == p2 * 2 || (p + p2) == p1 * 2)
    })
}

fn part2() -> Result<()> {
    run(|(p1, p2): &(Position, Position), p: &Position| -> bool {
        if p1 == p2 {
            false
        } else if p == p1 || p == p2 {
            true
        } else {
            let diff = p - p1;
            let d = p1 - p2;
            diff.0 % d.0 == 0 && diff.1 % d.1 == 0 && diff.0 / d.0 == diff.1 / d.1
        }
    })
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
