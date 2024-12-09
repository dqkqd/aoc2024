use std::{
    collections::{BTreeMap, BTreeSet, BinaryHeap},
    io::Read,
};

use anyhow::Result;

use rust::read;

fn part1() -> Result<()> {
    let mut reader = read(9, false)?;
    let mut s = String::new();
    reader.read_to_string(&mut s)?;
    let s = s.trim();

    let mut location = 0;
    let mut nums = vec![];
    let mut dots = vec![];
    for (i, x) in s.chars().enumerate() {
        let n = x.to_string().parse::<usize>()?;
        if n == 0 {
            continue;
        }
        if i % 2 == 0 {
            for _ in 0..n {
                nums.push((location, i / 2));
                location += 1;
            }
        } else {
            dots.extend(location..location + n);
            location += n;
        }
    }

    let mut dots = BinaryHeap::from_iter(dots.into_iter().map(|d| -(d as i32)));
    let mut moved_nums = vec![];

    while let Some((loc_num, num)) = nums.pop() {
        dots.push(-(loc_num as i32));

        let dot_loc = dots.pop().unwrap();
        let dot_loc = (-dot_loc) as usize;
        moved_nums.push((dot_loc, num));
        if loc_num == dot_loc {
            break;
        }
    }

    let ans = moved_nums
        .iter()
        .chain(nums.iter())
        .map(|(loc, num)| loc * num)
        .sum::<usize>();

    println!("{}", ans);

    Ok(())
}

fn part2() -> Result<()> {
    let mut reader = read(9, false)?;
    let mut s = String::new();
    reader.read_to_string(&mut s)?;
    let s = s.trim();
    // dbg!(s);

    let mut location = 0;
    let mut nums = vec![];
    let mut dots: BTreeMap<usize, BTreeSet<usize>> = BTreeMap::new();
    for (i, x) in s.chars().enumerate() {
        let n = x.to_string().parse::<usize>()?;
        if n == 0 {
            continue;
        }
        if i % 2 == 0 {
            nums.push((i / 2, n, location));
        } else {
            dots.entry(n).or_default().insert(location);
        }
        location += n;
    }

    let calculate = |num, size, loc_from| {
        let acc = loc_from * size + (size - 1) * size / 2;
        acc * num
    };

    let mut ans = 0;

    for (num, num_size, num_loc_from) in nums.into_iter().rev() {
        match dots
            .range_mut(num_size..)
            .filter_map(|(dot_size, dot_locs)| {
                dot_locs
                    .range(..num_loc_from)
                    .next()
                    .map(|dot_loc| (dot_size, dot_loc))
            })
            .min_by_key(|(_, dot_loc)| *dot_loc)
        {
            Some((&dot_size, _)) => {
                let dot_loc = dots.entry(dot_size).or_default().pop_first().unwrap();
                ans += calculate(num, num_size, dot_loc);
                dots.entry(dot_size - num_size)
                    .or_default()
                    .insert(dot_loc + num_size);
            }
            None => {
                ans += calculate(num, num_size, num_loc_from);
            }
        }
    }

    println!("{}", ans);

    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
