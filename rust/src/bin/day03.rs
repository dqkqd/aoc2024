use std::{
    fs::File,
    io::{BufReader, Read, Seek},
};

use anyhow::Result;

use anyhow::bail;
use rust::read;

#[derive(Debug)]
enum Instruction {
    Do,
    Dont,
    Invalid,
    Mul { lhs: i32, rhs: i32 },
}

fn read_one(reader: &mut BufReader<File>) -> Result<u8> {
    let mut buf = [0];
    reader.read_exact(&mut buf)?;
    Ok(buf[0])
}

fn read_while(reader: &mut BufReader<File>, pred: fn(&u8) -> bool) -> Vec<u8> {
    let mut out = vec![];

    while let Ok(c) = read_one(reader) {
        if !pred(&c) {
            reader.seek_relative(-1).expect("This must pass");
            break;
        }
        out.push(c);
    }

    out
}

fn expect_string(reader: &mut BufReader<File>, s: &str) -> anyhow::Result<()> {
    let pos = reader.stream_position()?;

    let expected = s.as_bytes();
    let mut buf = vec![0u8; expected.len()];

    reader.read_exact(&mut buf)?;

    if buf != expected {
        reader.seek(std::io::SeekFrom::Start(pos))?;
        bail!("Unexpected")
    }

    Ok(())
}

fn read_mul_instruction(reader: &mut BufReader<File>) -> Result<Instruction> {
    read_while(reader, |b| b != &b'm');
    expect_string(reader, "mul(")?;

    let lhs = read_while(reader, |b| (&b'0'..=&b'9').contains(&b));
    let lhs = String::from_utf8(lhs)?.parse::<i32>()?;

    expect_string(reader, ",")?;

    let rhs = read_while(reader, |b| (&b'0'..=&b'9').contains(&b));
    let rhs = String::from_utf8(rhs)?.parse::<i32>()?;

    expect_string(reader, ")")?;

    Ok(Instruction::Mul { lhs, rhs })
}

fn read_instruction(reader: &mut BufReader<File>) -> Result<Instruction> {
    read_while(reader, |b| b != &b'm' && b != &b'd');

    if expect_string(reader, "do()").is_ok() {
        return Ok(Instruction::Do);
    }
    if expect_string(reader, "don't()").is_ok() {
        return Ok(Instruction::Dont);
    }

    match read_mul_instruction(reader) {
        Ok(mul) => Ok(mul),
        Err(_) => {
            read_one(reader)?;
            Ok(Instruction::Invalid)
        }
    }
}

fn read_instructions() -> Result<Vec<Instruction>> {
    let mut reader = read(3, false)?;

    let mut instructions = vec![];
    while let Ok(ins) = read_instruction(&mut reader) {
        instructions.push(ins)
    }
    Ok(instructions)
}

fn part1() -> Result<()> {
    let instructions = read_instructions()?;
    let sum = instructions
        .into_iter()
        .map(|ins| match ins {
            Instruction::Mul { lhs, rhs } => lhs * rhs,
            _ => 0,
        })
        .sum::<i32>();
    println!("{}", sum);
    Ok(())
}

fn part2() -> Result<()> {
    let instructions = read_instructions()?;
    let mut sum = 0;

    let mut last = Instruction::Do;
    for ins in instructions {
        match ins {
            ins @ (Instruction::Do | Instruction::Dont) => last = ins,
            Instruction::Mul { lhs, rhs } => {
                if matches!(last, Instruction::Do) {
                    sum += lhs * rhs
                }
            }
            _ => {}
        }
    }
    println!("{}", sum);
    Ok(())
}

fn main() -> Result<()> {
    part1()?;
    part2()?;
    Ok(())
}
