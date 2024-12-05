use std::{
    collections::{BTreeMap, BTreeSet},
    io::BufRead,
    num::ParseIntError,
    str::FromStr,
};

use anyhow::Result;

use rust::read;

#[derive(Default, Debug)]
struct Rules {
    adj: BTreeMap<usize, BTreeSet<usize>>,
}

#[derive(Debug)]
struct Page {
    data: Vec<usize>,
}

impl Rules {
    fn add_rule(&mut self, x: usize, y: usize) {
        self.adj.entry(x).or_default().insert(y);
    }

    fn _dfs(
        &self,
        v: usize,
        visited: &mut BTreeSet<usize>,
        answer: &mut Vec<usize>,
        subset: &BTreeSet<usize>,
    ) {
        visited.insert(v);
        if let Some(adj) = self.adj.get(&v) {
            for u in adj.intersection(subset) {
                if !visited.contains(u) {
                    self._dfs(*u, visited, answer, subset);
                }
            }
        }
        answer.push(v);
    }

    fn topological_sort(&self, subset: BTreeSet<usize>) -> BTreeMap<usize, usize> {
        let mut visited: BTreeSet<usize> = BTreeSet::new();
        let mut answer = Vec::new();

        for u in &subset {
            if !visited.contains(u) {
                self._dfs(*u, &mut visited, &mut answer, &subset);
            }
        }
        answer.reverse();

        let mut map = BTreeMap::new();
        for (index, v) in answer.iter().enumerate() {
            map.insert(*v, index);
        }
        map
    }
}

impl FromStr for Page {
    type Err = ParseIntError;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        let data = s
            .split(",")
            .filter_map(|s| s.parse::<usize>().ok())
            .collect();
        Ok(Page { data })
    }
}

impl Page {
    fn good(&self, rules: &Rules) -> Option<usize> {
        let subset: BTreeSet<usize> = self.data.iter().cloned().collect();
        let sorted = rules.topological_sort(subset);
        let order: Vec<usize> = self.data.iter().map(|v| sorted[v]).collect();
        if order.windows(2).all(|v| v[0] <= v[1]) {
            Some(self.data[self.data.len() / 2])
        } else {
            None
        }
    }

    fn bad(&self, rules: &Rules) -> Option<usize> {
        let subset: BTreeSet<usize> = self.data.iter().cloned().collect();
        let sorted = rules.topological_sort(subset);
        let order: Vec<usize> = self.data.iter().map(|v| sorted[v]).collect();
        if order.windows(2).all(|v| v[0] <= v[1]) {
            None
        } else {
            let mut order = order;
            order.sort();
            let mid = order[order.len() / 2];
            let (key, _) = sorted.iter().find(|&(_, value)| value == &mid)?;
            Some(*key)
        }
    }
}

fn read_graph(sample: bool) -> Result<(Rules, Vec<Page>)> {
    let mut rules = Rules::default();

    for line in read(5, sample)?
        .lines()
        .take_while(|line| line.as_ref().map(|l| !l.trim().is_empty()).unwrap_or(false))
    {
        let line: Vec<usize> = line?
            .split("|")
            .map(|s| s.parse::<usize>().unwrap())
            .take(2)
            .collect();
        rules.add_rule(line[0], line[1]);
    }

    let mut pages = vec![];

    for line in read(5, sample)?
        .lines()
        .skip_while(|line| line.as_ref().map(|l| !l.trim().is_empty()).unwrap_or(false))
    {
        let line = line?;
        if line.trim().is_empty() {
            continue;
        }
        let page = Page::from_str(&line)?;
        pages.push(page);
    }

    Ok((rules, pages))
}

fn part1() -> Result<()> {
    let (rules, pages) = read_graph(false)?;
    let cnt = pages.iter().filter_map(|p| p.good(&rules)).sum::<usize>();
    println!("{:?}", cnt);
    Ok(())
}

fn part2() -> Result<()> {
    let (rules, pages) = read_graph(false)?;
    let cnt = pages.iter().filter_map(|p| p.bad(&rules)).sum::<usize>();
    println!("{:?}", cnt);
    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
