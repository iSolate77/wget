use anyhow::Result;
use clap::Parser;
use reqwest::{Client, Error, Response};
use std::env;

#[derive(Parser, Debug)]
pub struct Args {
    /// Download URL in the background
    #[clap(long, short, default_value = "wget-log")]
    background: String,
}

pub async fn get_url(url: &str) -> Result<Response, Error> {
    Client::new().get(url).send().await
}

pub fn parse_args() -> Result<Args> {
    let args = Args::parse();
    Ok(args)
}
