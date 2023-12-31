// mod lib;
use anyhow::Result;
use clap::Parser;
use futures_util::StreamExt;
use indicatif::{ProgressBar, ProgressState, ProgressStyle};
use reqwest::{Client, Error, Response, ResponseBuilderExt, Url};
use std::fs::File;
use std::io::Write;
use tokio::io::AsyncWriteExt;

#[derive(Parser, Debug)]
struct Args {
    /// URL
    #[clap(required = true)]
    url: String,

    /// Go to background immediately after startup. If no output file is specified, the default is to "wget-log".
    #[clap(short = 'B', long)]
    background: bool,

    /// Specify the output file name. The default is to write to the console.
    #[clap(short = 'O', long)]
    output: Option<String>,

    /// Specify a path to store the downloaded file
    #[clap(short = 'P', long, default_value = ".")]
    path: String,

    /// Specify the rate limit for the download. The rate is measured in bytes / sec. The default is 0, which means no limit.
    #[clap(long)]
    rate_limit: bool,

    /// Specify an input file containing URLs to download. If you specify "-", the URL is read from the standard input.
    #[clap(short = 'i', long)]
    input_file: bool,

    /// Specify a mirror site for the download. The default is no mirror.
    #[clap(long)]
    mirror: bool,

    /// Reject a certain type of file.
    #[clap(short = 'R', requires = "mirror")]
    reject: Option<String>,

    /// Exclude a certain type of file.
    #[clap(short = 'X', requires = "mirror")]
    exclude: Option<String>,
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();
    let mut final_url = args.url.clone();
    if !final_url.starts_with("http://") && !final_url.starts_with("https://") {
        final_url = format!("http://{}", final_url);
    }
    if args.mirror {
        mirror_website(&final_url.clone()).await?;
    }
    let output_file_name = args.output.unwrap_or_else(|| {
        derive_file_name_from_url(&final_url).unwrap_or_else(|| "index.html".to_string())
    });

    download_url(&final_url, Some(&output_file_name)).await?;

    Ok(())
}

async fn download_url(url: &str, output_file: Option<&str>) -> Result<()> {
    let response = Client::new().get(url).send().await?;

    let content_size = response
        .headers()
        .get(reqwest::header::CONTENT_LENGTH)
        .and_then(|value| value.to_str().ok())
        .and_then(|value| value.parse().ok());

    let total_size = content_size.unwrap_or(0);

    let pb = ProgressBar::new(total_size);
    pb.set_style(ProgressStyle::with_template("{spinner:.green} [{elapsed_precise}] [{wide_bar:.cyan/blue}] {bytes}/{total_bytes} ({eta})")
        .unwrap()
        // .with_key("eta", |state: &ProgressState, w: &mut dyn Write| write!(w, "{:.1}s", state.eta().as_secs_f64()).unwrap())
        .progress_chars("#>-"));

    let mut stream = response.bytes_stream();

    match output_file {
        Some(file_name) => {
            let mut file = File::create(file_name)?;
            while let Some(chunk) = stream.next().await {
                let data = chunk?;
                file.write_all(&data)?;
                pb.inc(data.len() as u64);
            }
        }
        None => {
            let mut content = Vec::new();
            while let Some(chunk) = stream.next().await {
                let data = chunk?;
                content.extend_from_slice(&data);
                pb.inc(data.len() as u64);
            }
        }
    }

    pb.finish_with_message("Download complete");
    Ok(())
}

fn save_to_file(file_name: &str, content: &str) -> Result<(), std::io::Error> {
    let mut file = File::create(file_name)?;
    file.write_all(content.as_bytes())?;
    Ok(())
}

fn derive_file_name_from_url(url: &str) -> Option<String> {
    Url::parse(url)
        .ok()?
        .path_segments()?
        .last()
        .and_then(|last_segment| {
            if last_segment.contains('.') {
                Some(last_segment.to_string())
            } else {
                None
            }
        })
}

async fn mirror_website(url: &str) -> Result<()> {
    let client = Client::new();
    let base_url = Url::parse(url)?;
    let host_name = base_url.host_str().unwrap_or("website");

    let base_path = Path::new(host_name);
    if !base_path.exists() {
        fs::create_dir_all(base_path)?;
    }

    recursive_download(&client, &base_url, base_path).await?;

    Ok(())
}

async fn recursive_download(
    client: &Client,
    url: &str,
    visited: &mut HashSet<String>,
) -> Result<()> {
    if visited.contains(url) {
        return Ok(());
    }

    // Download the page content
    let content = download_text(client, url).await?;
    visited.insert(url.to_string());

    let html = Html::parse_document(&content);
    let base_url = Url::parse(url)?;

    // Find and process all links and images
    process_links_and_images(client, &html, &base_url, visited).await?;

    Ok(())
}

async fn process_links_and_images(
    client: &Client,
    html: &Html,
    base_url: &Url,
    visited: &mut HashSet<String>,
) -> Result<()> {
    let link_selector = Selector::parse("a[href]").unwrap();
    let image_selector = Selector::parse("img[src]").unwrap();

    for element in html.select(&link_selector) {
        if let Some(href) = element.value().attr("href") {
            let next_url = base_url.join(href)?;
            recursive_download(client, next_url.as_str(), visited).await?;
        }
    }

    for element in html.select(&image_selector) {
        if let Some(src) = element.value().attr("src") {
            let image_url = base_url.join(src)?;
            download_image(client, image_url.as_str()).await?;
        }
    }

    Ok(())
}

async fn download_image(client: &Client, url: &str) -> Result<()> {
    let res = client.get(url).send().await?;
    let content = res.bytes().await?;
    save_bytes_to_file(url, &content)?;
    Ok(())
}


fn save_bytes_to_file(file_name: &str, content: &[u8]) -> Result<(), std::io::Error> {
    let mut file = File::create(file_name)?;
    file.write_all(content)?;
    Ok(())
}
