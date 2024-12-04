use std::io::BufRead;

use anyhow::Result;

use rust::read;

fn read_map() -> Result<Vec<Vec<char>>> {
    let reader = read(4, false)?;
    let map = reader
        .lines()
        .map(|f| f.unwrap().chars().collect())
        .collect();
    Ok(map)
}

fn xmas_count(map: &[Vec<char>], x: i32, y: i32) -> i32 {
    let mut cnt = 0;

    let h = map.len() as i32;
    let w = map[0].len() as i32;

    for dx in -1..=1 {
        for dy in -1..=1 {
            let mut pos = vec![];
            for i in 0..4 {
                pos.push((x + i * dx, y + i * dy));
            }
            if pos
                .iter()
                .all(|&(px, py)| px >= 0 && px < h && py >= 0 && py < w)
            {
                let mut s = String::new();
                for (px, py) in pos {
                    let c = map[px as usize][py as usize];
                    s.push(c);
                }

                if s == "XMAS" {
                    cnt += 1;
                }
            }
        }
    }
    cnt
}

fn xmas_3x3_count(map: &[Vec<char>], x: usize, y: usize) -> i32 {
    let mut s = vec![];
    let h = map.len();
    let w = map[0].len();
    if x + 3 > h {
        return 0;
    }
    if y + 3 > w {
        return 0;
    }

    for v in map.iter().skip(x).take(3) {
        let v = v.get(y..y + 3).unwrap();
        s.push(v);
    }

    if s[1][1] != 'A' {
        return 0;
    }

    let pos = [
        (0, 0),
        (0, 2),
        (2, 2),
        (2, 0),
        (0, 0),
        (0, 2),
        (2, 2),
        (2, 0),
    ];
    for offset in 0..4 {
        let mut check = String::new();
        for &(px, py) in &pos[offset..offset + 4] {
            check.push(s[px][py]);
        }
        if check == "MMSS" {
            return 1;
        }
    }

    0
}

fn part1() -> Result<()> {
    let map = read_map()?;
    let h = map.len();
    let w = map[0].len();

    let mut cnt = 0;

    for x in 0..h {
        for y in 0..w {
            cnt += xmas_count(&map, x as i32, y as i32);
        }
    }
    println!("{:?}", cnt);

    Ok(())
}

fn part2() -> Result<()> {
    let map = read_map()?;
    let h = map.len();
    let w = map[0].len();

    let mut cnt = 0;

    for x in 0..h {
        for y in 0..w {
            cnt += xmas_3x3_count(&map, x, y);
        }
    }
    println!("{:?}", cnt);

    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
