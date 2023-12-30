use reqwest::{Client, Error, Response};

fn get_url(url: &str) -> Result<Response, Error> {
    Client::new().get(url).send().await
}
