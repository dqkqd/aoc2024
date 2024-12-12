use std::{
    collections::{BTreeMap, BTreeSet, VecDeque},
    io::BufRead,
};

use anyhow::Result;

use rust::read;

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord)]
enum Edge {
    Vertical(usize, usize),
    Horizontal(usize, usize),
}

#[derive(Debug)]
struct Map {
    h: usize,
    w: usize,
    data: Vec<Vec<char>>,
}

fn adjacent_edges((x, y): (usize, usize)) -> [Edge; 4] {
    [
        Edge::Horizontal(x, y),
        Edge::Horizontal(x + 1, y),
        Edge::Vertical(x, y),
        Edge::Vertical(x, y + 1),
    ]
}

impl Map {
    fn good_point(&self, (x, y): &(i32, i32)) -> bool {
        *x >= 0 && *x < self.h as i32 && *y >= 0 && *y < self.w as i32
    }

    fn adjacent_points(&self, (x, y): (usize, usize)) -> Vec<(usize, usize)> {
        let x = x as i32;
        let y = y as i32;
        [(x, y - 1), (x, y + 1), (x - 1, y), (x + 1, y)]
            .into_iter()
            .filter(|pos| self.good_point(pos))
            .map(|(x, y)| (x as usize, y as usize))
            .collect()
    }

    fn is_separated_edge(&self, edge: &Edge) -> bool {
        let (lhs, rhs) = match edge {
            Edge::Vertical(x, y) => {
                let x = *x as i32;
                let y = *y as i32;

                let lhs = (x, y - 1);
                let rhs = (x, y);

                (lhs, rhs)
            }
            Edge::Horizontal(x, y) => {
                let x = *x as i32;
                let y = *y as i32;

                let lhs = (x - 1, y);
                let rhs = (x, y);

                (lhs, rhs)
            }
        };

        if self.good_point(&lhs) && self.good_point(&rhs) {
            self.data[lhs.0 as usize][lhs.1 as usize] != self.data[rhs.0 as usize][rhs.1 as usize]
        } else {
            true
        }
    }

    fn visit(
        &self,
        pos: (usize, usize),
        colors: &mut [Vec<Option<i32>>],
    ) -> (BTreeSet<(usize, usize)>, BTreeSet<Edge>) {
        let datum = self.data[pos.0][pos.1];
        let color = colors[pos.0][pos.1];

        let mut visited = BTreeSet::new();
        let mut edges = BTreeSet::new();

        let mut queue = VecDeque::from([pos]);
        while let Some(pos) = queue.pop_front() {
            if visited.contains(&pos) {
                continue;
            }
            visited.insert(pos);
            edges.extend(
                adjacent_edges(pos)
                    .into_iter()
                    .filter(|e| self.is_separated_edge(e)),
            );

            let adjacents: Vec<(usize, usize)> = self
                .adjacent_points(pos)
                .into_iter()
                .filter(|(x, y)| self.data[*x][*y] == datum)
                .collect();
            for (x, y) in &adjacents {
                colors[*x][*y] = color;
            }
            queue.extend(adjacents);
        }

        (visited, edges)
    }

    fn price(&self, pos: (usize, usize), colors: &mut [Vec<Option<i32>>]) -> usize {
        let (visited, edges) = self.visit(pos, colors);
        let area = visited.len();
        let perimeter = edges
            .into_iter()
            .map(|edge| if self.is_separated_edge(&edge) { 1 } else { 0 })
            .sum::<usize>();

        area * perimeter
    }

    fn prices(self) -> usize {
        let mut colors: Vec<Vec<Option<i32>>> =
            vec![vec![None; self.data[0].len()]; self.data.len()];

        let mut ans = 0;
        let mut color = 0;
        for x in 0..self.data.len() {
            for y in 0..self.data[0].len() {
                if colors[x][y].is_some() {
                    continue;
                }
                colors[x][y] = Some(color);
                let price = self.price((x, y), &mut colors);
                ans += price;
                color += 1;
            }
        }

        ans
    }

    fn edge_same_side(&self, lhs: &Edge, rhs: &Edge) -> bool {
        match (lhs.clone(), rhs.clone()) {
            (Edge::Vertical(x1, y1), Edge::Vertical(x2, y2))
                if y1 == y2 && (x1 == x2 + 1 || x2 == x1 + 1) =>
            {
                (y1 > 0 && self.data[x1][y1 - 1] == self.data[x2][y2 - 1])
                    || (y1 < self.w && self.data[x1][y1] == self.data[x2][y2])
            }
            (Edge::Horizontal(x1, y1), Edge::Horizontal(x2, y2))
                if x1 == x2 && (y1 == y2 + 1 || y2 == y1 + 1) =>
            {
                (x1 > 0 && self.data[x1 - 1][y1] == self.data[x2 - 1][y2])
                    || (x1 < self.h && self.data[x1][y1] == self.data[x2][y2])
            }
            _ => false,
        }
    }

    fn price2(&self, pos: (usize, usize), colors: &mut [Vec<Option<i32>>]) -> usize {
        let (visited, edges) = self.visit(pos, colors);
        let area = visited.len();

        let mut groups: BTreeMap<usize, Vec<Edge>> = BTreeMap::new();
        let mut sides = 0;
        for edge in edges.into_iter() {
            if let Some((_, g)) = groups
                .iter_mut()
                .find(|(_, g)| g.iter().any(|e| self.edge_same_side(&edge, e)))
            {
                g.push(edge)
            } else {
                groups.insert(sides, vec![edge]);
                sides += 1;
            }
        }

        area * sides
    }

    fn prices2(self) -> usize {
        let mut colors: Vec<Vec<Option<i32>>> =
            vec![vec![None; self.data[0].len()]; self.data.len()];

        let mut ans = 0;
        let mut color = 0;
        for x in 0..self.data.len() {
            for y in 0..self.data[0].len() {
                if colors[x][y].is_some() {
                    continue;
                }
                colors[x][y] = Some(color);
                let price = self.price2((x, y), &mut colors);
                ans += price;
                color += 1;
            }
        }

        ans
    }
}

fn read_map() -> Result<Map> {
    let reader = read(12, false)?;
    let data: Vec<Vec<char>> = reader
        .lines()
        .map_while(Result::ok)
        .filter(|line| !line.is_empty())
        .map(|line| line.chars().collect())
        .collect();
    Ok(Map {
        h: data.len(),
        w: data[0].len(),
        data,
    })
}

fn part1() -> Result<()> {
    let map = read_map()?;
    println!("{}", map.prices());
    Ok(())
}

fn part2() -> Result<()> {
    let map = read_map()?;
    println!("{}", map.prices2());
    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
