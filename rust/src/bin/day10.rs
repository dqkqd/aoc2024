use std::{
    collections::{BTreeSet, VecDeque},
    io::BufRead,
};

use anyhow::Result;

use rust::read;

#[derive(Debug)]
struct Map {
    data: Vec<Vec<u8>>,
}

impl Map {
    fn zeros(&self) -> Vec<(usize, usize)> {
        self.data
            .iter()
            .enumerate()
            .flat_map(|(x, line)| {
                line.iter()
                    .enumerate()
                    .filter(|(_, c)| c == &&0)
                    .map(|(y, _)| (x, y))
                    .collect::<Vec<(usize, usize)>>()
            })
            .collect()
    }
    fn adj_points(&self, (x, y): (usize, usize)) -> Vec<(usize, usize)> {
        let (x, y) = (x as i32, y as i32);
        let pos = [(x + 1, y), (x, y + 1), (x - 1, y), (x, y - 1)];

        let good = |(x, y): (&i32, &i32)| {
            *x >= 0 && *x < self.data.len() as i32 && *y >= 0 && *y < self.data[0].len() as i32
        };

        pos.into_iter()
            .filter(|(x, y)| good((x, y)))
            .map(|(x, y)| (x as usize, y as usize))
            .collect()
    }

    fn count_score(&self, pos: (usize, usize)) -> usize {
        let mut visited: BTreeSet<(usize, usize)> = BTreeSet::new();
        let mut queue = VecDeque::new();

        queue.push_back(pos);
        while let Some(pos) = queue.pop_front() {
            if visited.contains(&pos) {
                continue;
            }
            visited.insert(pos);
            let value = self.data[pos.0][pos.1];
            let good_points = self
                .adj_points(pos)
                .into_iter()
                .filter(|p| !visited.contains(p))
                .filter(|(x, y)| value + 1 == self.data[*x][*y]);
            queue.extend(good_points);
        }

        visited
            .into_iter()
            .filter(|(x, y)| self.data[*x][*y] == 9)
            .count()
    }

    fn count_score_part2(&self, pos: (usize, usize)) -> usize {
        let mut visited: BTreeSet<(usize, usize)> = BTreeSet::new();
        let mut score_map = vec![vec![0; self.data[0].len()]; self.data.len()];
        let mut queue = VecDeque::new();

        score_map[pos.0][pos.1] = 1;
        queue.push_back(pos);
        while let Some(pos) = queue.pop_front() {
            if visited.contains(&pos) {
                continue;
            }
            visited.insert(pos);
            let value = self.data[pos.0][pos.1];
            let score = score_map[pos.0][pos.1];
            self.adj_points(pos)
                .into_iter()
                .filter(|p| !visited.contains(p))
                .filter(|(x, y)| value + 1 == self.data[*x][*y])
                .for_each(|(x, y)| {
                    score_map[x][y] += score;
                    queue.push_back((x, y));
                });
        }
        visited
            .into_iter()
            .filter(|(x, y)| self.data[*x][*y] == 9)
            .map(|(x, y)| score_map[x][y])
            .sum()
    }
}

fn read_map() -> Result<Map> {
    let reader = read(10, false)?;
    let data = reader
        .lines()
        .map_while(Result::ok)
        .filter(|s| !s.is_empty())
        .map(|s| s.split("").filter_map(|c| c.parse::<u8>().ok()).collect())
        .collect();
    Ok(Map { data })
}

fn part1() -> Result<()> {
    let map = read_map()?;
    let zeros = map.zeros();
    let ans: usize = zeros.into_iter().map(|p| map.count_score(p)).sum();
    println!("{}", ans);
    Ok(())
}

fn part2() -> Result<()> {
    let map = read_map()?;
    let zeros = map.zeros();
    let ans: usize = zeros.into_iter().map(|p| map.count_score_part2(p)).sum();
    println!("{}", ans);
    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
