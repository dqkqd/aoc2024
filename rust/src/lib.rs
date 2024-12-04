use anyhow::Result;
use std::{env::current_dir, fs::File, io::BufReader};

pub fn read(day: u16, sample: bool) -> Result<BufReader<File>> {
    let input_folder = current_dir()?.parent().unwrap().join("input");
    let input_filename = if sample { "sample.txt" } else { "input.txt" };
    let file = input_folder
        .join(format!("day{:0>2}", day))
        .join(input_filename);

    let f = File::open(file)?;
    Ok(BufReader::new(f))
}
