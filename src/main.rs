// mod lib;
use std::env;
use std::io::{self, Write};

const HELP_MESSAGE: &str = "reboot wget 0.1.0\nUsage: wget [OPTION]... [URL]...\nMandatory arguments to long options are mandatory for short options too.";
const NO_OPTIONS: &str =
    "wget: missing URL\nUsage: wget [OPTION]... [URL]...\nTry `wget --help' for more options.";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<_> = env::args().skip(1).collect();
    let mut lock = io::stdout().lock();
    if args.len() < 1 {
        writeln!(lock, "{}", NO_OPTIONS)?;
    }
    if args[0] == "--help" || args[0] == "-h" {
        writeln!(lock, "{}", HELP_MESSAGE)?;
    }
    let options: Vec<_> = args
        .iter()
        .filter_map(|x| {
            if x.starts_with("--") || x.starts_with("-") {
                Some(x)
            } else {
                None
            }
        })
        .collect();
    println!("{:?} {:?}", options, args);
    Ok(())
}
