use std::{
    fs::File,
    io::{BufRead, BufReader},
};

use anyhow::{bail, Result};

use rust::read;

fn extended_euclid(a: i64, b: i64) -> (i64, i64, i64) {
    debug_assert!(a >= 0 && b >= 0);
    if b == 0 {
        (a, 1, 0)
    } else {
        let (g, x, y) = extended_euclid(b, a % b);
        (g, y, x - y * (a / b))
    }
}

#[derive(Debug)]
enum Root {
    No,
    Yes(i64, i64),
    Any,
}

#[derive(Debug, PartialEq, Eq, PartialOrd, Ord)]
struct LinearEquation1 {
    a: i64,
    b: i64,
}

#[derive(Debug, PartialEq, Eq, PartialOrd, Ord)]
struct LinearEquation2 {
    a: i64,
    b: i64,
    c: i64,
}

#[derive(Debug)]
struct LinearSystem {
    eq1: LinearEquation2,
    eq2: LinearEquation2,
}

impl LinearEquation1 {
    fn solve(&self) -> Root {
        if self.a == 0 && self.b == 0 {
            Root::Any
        } else if self.a != 0 && self.b % self.a == 0 {
            Root::Yes(0, self.b / self.a)
        } else {
            Root::No
        }
    }
}

impl LinearEquation2 {
    fn solve(&self) -> (Root, Root) {
        if self.a == 0 {
            let eq = LinearEquation1 {
                a: self.b,
                b: self.c,
            };
            (Root::Any, eq.solve())
        } else if self.b == 0 {
            let eq = LinearEquation1 {
                a: self.a,
                b: self.c,
            };
            (eq.solve(), Root::Any)
        } else {
            let (g, xg, yg) = extended_euclid(self.a, self.b);
            debug_assert!(self.a * xg + self.b * yg == g);

            if self.c % g != 0 {
                (Root::No, Root::No)
            } else {
                let x0 = xg * self.c / g;
                let y0 = yg * self.c / g;
                (Root::Yes(self.b / g, x0), Root::Yes(-self.a / g, y0))
            }
        }
    }
}

impl LinearSystem {
    fn solve(&self) -> (Root, Root) {
        let (x, y) = self.eq1.solve();
        match (x, y) {
            (Root::No, _) | (_, Root::No) => (Root::No, Root::No),

            (Root::Any, Root::Any) => self.eq2.solve(),

            (x @ Root::Yes(x0, x1), y @ Root::Yes(y0, y1)) => {
                // a2 * (x0 * t + x1) + b2 * (y0 * t + y1) == c2
                let eq = LinearEquation1 {
                    a: self.eq2.a * x0 + self.eq2.b * y0,
                    b: self.eq2.c - (self.eq2.a * x1 + self.eq2.b * y1),
                };

                match eq.solve() {
                    Root::No => (Root::No, Root::No),
                    Root::Yes(t0, t1) => {
                        debug_assert!(t0 == 0);
                        (Root::Yes(0, t1 * x0 + x1), Root::Yes(0, t1 * y0 + y1))
                    }
                    Root::Any => (x, y),
                }
            }
            (x @ Root::Yes(x0, x1), Root::Any) => {
                debug_assert!(self.eq1.b == 0 && x0 == 0);

                // a2 * x1 + b2 * y = c2
                let eq = LinearEquation1 {
                    a: self.eq2.b,
                    b: self.eq2.c - self.eq2.a * x1,
                };
                let y = eq.solve();
                match y {
                    Root::No => (Root::No, Root::No),
                    y => (x, y),
                }
            }
            (Root::Any, y @ Root::Yes(y0, y1)) => {
                debug_assert!(self.eq1.a == 0 && y0 == 0);

                // a2 * x + b2 * y1 = c2
                let eq = LinearEquation1 {
                    a: self.eq2.a,
                    b: self.eq2.c - self.eq2.b * y1,
                };
                let x = eq.solve();
                match x {
                    Root::No => (Root::No, Root::No),
                    x => (x, y),
                }
            }
        }
    }
}

fn read_linear_system(reader: &mut BufReader<File>) -> Result<LinearSystem> {
    let mut line1 = String::new();
    let mut line2 = String::new();
    let mut line3 = String::new();
    reader.read_line(&mut line1)?;
    reader.read_line(&mut line2)?;
    reader.read_line(&mut line3)?;

    if line1.is_empty() || line2.is_empty() || line3.is_empty() {
        bail!("stop")
    }

    let get_num = |line: &str, sep: &str| -> (i64, i64) {
        let values: Vec<i64> = line
            .split(":")
            .last()
            .unwrap()
            .split(",")
            .take(2)
            .map(|s| s.split(sep).last().unwrap().trim().parse::<i64>().unwrap())
            .collect();
        (values[0], values[1])
    };

    let (a1, a2) = get_num(&line1, "+");
    let (b1, b2) = get_num(&line2, "+");
    let (c1, c2) = get_num(&line3, "=");

    // read empty line
    reader.read_line(&mut line3)?;

    Ok(LinearSystem {
        eq1: LinearEquation2 {
            a: a1,
            b: b1,
            c: c1,
        },
        eq2: LinearEquation2 {
            a: a2,
            b: b2,
            c: c2,
        },
    })
}

fn part1() -> Result<()> {
    let mut reader = read(13, false)?;
    let mut ans = 0;
    while let Ok(system) = read_linear_system(&mut reader) {
        let (x, y) = system.solve();
        match (x, y) {
            (Root::Yes(x0, x1), Root::Yes(y0, y1)) => {
                debug_assert!(x0 == 0 && y0 == 0);
                if (0..=100).contains(&x1) && (0..=100).contains(&y1) {
                    ans += 3 * x1 + y1;
                }
            }
            (_, Root::No) | (Root::No, _) | (Root::Any, Root::Any) => continue,
            (Root::Yes(_, _), Root::Any) => todo!(),
            (Root::Any, Root::Yes(_, _)) => todo!(),
        }
    }
    println!("{}", ans);
    Ok(())
}

fn part2() -> Result<()> {
    let mut reader = read(13, false)?;
    let mut ans = 0;
    while let Ok(system) = read_linear_system(&mut reader) {
        let system = LinearSystem {
            eq1: LinearEquation2 {
                a: system.eq1.a,
                b: system.eq1.b,
                c: system.eq1.c + 10000000000000,
            },
            eq2: LinearEquation2 {
                a: system.eq2.a,
                b: system.eq2.b,
                c: system.eq2.c + 10000000000000,
            },
        };
        let (x, y) = system.solve();
        match (x, y) {
            (Root::Yes(x0, x1), Root::Yes(y0, y1)) => {
                debug_assert!(x0 == 0 && y0 == 0);
                ans += 3 * x1 + y1;
            }
            (_, Root::No) | (Root::No, _) | (Root::Any, Root::Any) => continue,
            (Root::Yes(_, _), Root::Any) => todo!(),
            (Root::Any, Root::Yes(_, _)) => todo!(),
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
