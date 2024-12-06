use std::{
    collections::{BTreeMap, BTreeSet},
    io::BufRead,
};

use anyhow::Result;

use rust::read;
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum Direction {
    Up,
    Down,
    Left,
    Right,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
struct Position {
    x: i32,
    y: i32,
    direction: Direction,
}

struct Map {
    data: Vec<Vec<char>>,
    horizontal_obstacles: BTreeMap<i32, Vec<i32>>,
    vertical_obstacles: BTreeMap<i32, Vec<i32>>,
}

impl Direction {
    fn turn(self) -> Direction {
        match self {
            Direction::Up => Direction::Right,
            Direction::Down => Direction::Left,
            Direction::Left => Direction::Up,
            Direction::Right => Direction::Down,
        }
    }
}

impl Position {
    fn next_pos(&self) -> Position {
        match self.direction {
            Direction::Up => Position {
                x: self.x - 1,
                y: self.y,
                direction: self.direction,
            },
            Direction::Down => Position {
                x: self.x + 1,
                y: self.y,
                direction: self.direction,
            },
            Direction::Left => Position {
                x: self.x,
                y: self.y - 1,
                direction: self.direction,
            },
            Direction::Right => Position {
                x: self.x,
                y: self.y + 1,
                direction: self.direction,
            },
        }
    }

    fn is_outside(&self, map: &Map) -> bool {
        self.x < 0
            || self.y < 0
            || self.x >= map.data.len() as i32
            || self.y >= map.data[0].len() as i32
    }

    fn is_valid(&self, map: &Map) -> bool {
        !self.is_outside(map) && map.data[self.x as usize][self.y as usize] != '#'
    }

    fn turn(&mut self) {
        self.direction = self.direction.turn();
    }

    fn farthest_pos(&self, map: &Map) -> Option<Position> {
        debug_assert!(self.is_valid(map));
        match self.direction {
            Direction::Up => {
                let obstacles = map.vertical_obstacles.get(&self.y)?;
                let index = obstacles.binary_search(&self.x).unwrap_err();
                let index = index.checked_sub(1)?;
                Some(Position {
                    x: obstacles[index] + 1,
                    y: self.y,
                    direction: self.direction,
                })
            }
            Direction::Down => {
                let obstacles = map.vertical_obstacles.get(&self.y)?;
                let index = obstacles.binary_search(&self.x).unwrap_err();
                if index >= obstacles.len() {
                    None
                } else {
                    Some(Position {
                        x: obstacles[index] - 1,
                        y: self.y,
                        direction: self.direction,
                    })
                }
            }
            Direction::Left => {
                let obstacles = map.horizontal_obstacles.get(&self.x)?;
                let index = obstacles.binary_search(&self.y).unwrap_err();
                let index = index.checked_sub(1)?;
                Some(Position {
                    x: self.x,
                    y: obstacles[index] + 1,
                    direction: self.direction,
                })
            }
            Direction::Right => {
                let obstacles = map.horizontal_obstacles.get(&self.x)?;
                let index = obstacles.binary_search(&self.y).unwrap_err();
                if index >= obstacles.len() {
                    None
                } else {
                    Some(Position {
                        x: self.x,
                        y: obstacles[index] - 1,
                        direction: self.direction,
                    })
                }
            }
        }
    }
}

impl Map {
    fn new(data: Vec<Vec<char>>) -> Map {
        let mut horizontal_obstacles: BTreeMap<i32, Vec<i32>> = BTreeMap::new();
        let mut vertical_obstacles: BTreeMap<i32, Vec<i32>> = BTreeMap::new();
        for (x, vy) in data.iter().enumerate() {
            for (y, c) in vy.iter().enumerate() {
                let x = x as i32;
                let y = y as i32;
                if c == &'#' {
                    vertical_obstacles.entry(y).or_default().push(x);
                    horizontal_obstacles.entry(x).or_default().push(y);
                }
            }
        }

        Map {
            data,
            horizontal_obstacles,
            vertical_obstacles,
        }
    }
    fn run(&self, start: &Position) -> Option<Vec<(i32, i32)>> {
        let mut pos = *start;

        let mut positions = BTreeSet::new();

        if !pos.is_valid(self) {
            return None;
        }

        while !pos.is_outside(self) {
            if positions.contains(&pos) {
                return None;
            }
            positions.insert(pos);

            let next_pos = pos.next_pos();
            if next_pos.is_outside(self) {
                break;
            }
            if next_pos.is_valid(self) {
                pos = next_pos;
            } else {
                pos.turn();
            }
        }

        Some(
            positions
                .into_iter()
                .map(|p| (p.x, p.y))
                .collect::<BTreeSet<_>>()
                .into_iter()
                .collect(),
        )
    }

    fn mark_obstacle(&mut self, x: usize, y: usize) {
        debug_assert!(self.data[x][y] != '#');
        self.data[x][y] = '#';
        let x = x as i32;
        let y = y as i32;

        let vertical_obstacles = self.vertical_obstacles.entry(y).or_default();
        let index = vertical_obstacles.partition_point(|&v| v <= x);
        vertical_obstacles.insert(index, x);

        let horizontal_obstacles = self.horizontal_obstacles.entry(x).or_default();
        let index = horizontal_obstacles.partition_point(|&v| v <= y);
        horizontal_obstacles.insert(index, y);
    }
    fn unmark_obstacle(&mut self, x: usize, y: usize) {
        debug_assert!(self.data[x][y] == '#');
        self.data[x][y] = '.';
        let x = x as i32;
        let y = y as i32;

        let pos = self.vertical_obstacles[&y]
            .iter()
            .position(|v| v == &x)
            .expect("position must exist");
        self.vertical_obstacles.get_mut(&y).unwrap().remove(pos);

        let pos = self.horizontal_obstacles[&x]
            .iter()
            .position(|v| v == &y)
            .expect("position must exist");
        self.horizontal_obstacles.get_mut(&x).unwrap().remove(pos);
    }

    fn is_loop(&self, start: &Position) -> bool {
        let mut traces: BTreeSet<Position> = BTreeSet::new();
        let mut pos = *start;
        loop {
            if traces.contains(&pos) {
                return true;
            }
            traces.insert(pos);
            match pos.farthest_pos(self) {
                Some(p) => {
                    pos = p;
                    pos.turn();
                }
                None => return false,
            }
        }
    }
}

fn read_map() -> Result<(Map, Position)> {
    let reader = read(6, false)?;
    let mut data = vec![];
    for line in reader.lines() {
        let line = line?;
        if line.trim().is_empty() {
            continue;
        }
        let chars: Vec<char> = line.trim().chars().collect();
        data.push(chars)
    }

    let mut start = Position {
        x: 0,
        y: 0,
        direction: Direction::Up,
    };
    for (x, line) in data.iter().enumerate() {
        if let Some(y) = line.iter().position(|p| p == &'^') {
            start.x = x as i32;
            start.y = y as i32;
            break;
        }
    }

    let map = Map::new(data);

    Ok((map, start))
}

fn part1() -> Result<()> {
    let (map, start) = read_map()?;
    let positions = map.run(&start).expect("Must not a loop");
    println!("{}", positions.len());
    Ok(())
}

fn part2() -> Result<()> {
    let (mut map, start) = read_map()?;
    let positions = map.run(&start).expect("Must not a loop");

    let mut ans = 0;
    for (x, y) in positions {
        if x == start.x && y == start.y {
            continue;
        }

        let x = x as usize;
        let y = y as usize;
        map.mark_obstacle(x, y);
        if map.is_loop(&start) {
            ans += 1;
        }
        map.unmark_obstacle(x, y);
    }

    println!("{}", ans);
    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;

    Ok(())
}
