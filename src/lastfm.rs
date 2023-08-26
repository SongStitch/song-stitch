use image::RgbImage;
use miette::{miette, IntoDiagnostic};
use serde::Deserialize;

#[derive(Deserialize, Clone)]
pub struct Artist {
    url: String,
    pub name: String,
    mbid: String,
}

#[derive(Deserialize, Clone)]
pub struct Album {
    pub name: String,
    pub artist: Artist,
    pub playcount: String,
    url: String,
    #[serde(rename = "image")]
    images: Vec<Image>,
}

impl Album {
    pub async fn get_image(&self, size: &str) -> Option<RgbImage> {
        let image = self.images.iter().find(|image| image.size == size).unwrap();
        if image.url == "" {
            return None;
        }
        let img_bytes = reqwest::get(&image.url).await.ok()?.bytes().await.ok()?;
        Some(image::load_from_memory(&img_bytes).unwrap().to_rgb8())
    }
}

#[derive(Deserialize, Clone)]
pub struct Image {
    #[serde(rename = "#text")]
    url: String,
    size: String,
}

#[derive(Deserialize)]
pub struct TopAlbums {
    album: Vec<Album>,
}

#[derive(Deserialize)]
pub struct LastFmResponse {
    topalbums: TopAlbums,
}

pub struct LastFm {
    endpoint: String,
    api_key: String,
}

impl LastFm {
    pub fn new() -> Self {
        Self {
            endpoint: std::env::var("LASTFM_ENDPOINT").unwrap(),
            api_key: std::env::var("LASTFM_API_KEY").unwrap(),
        }
    }

    pub async fn get_top_albums(
        &self,
        username: &str,
        period: &str,
        count: usize,
    ) -> miette::Result<Vec<Album>> {
        println!("{} {} {}", username, period, count);
        let url = format!(
            "{}/?method=user.gettopalbums&user={}&period={}&api_key={}&limit={}&format=json",
            self.endpoint, username, period, self.api_key, count
        );
        let response = reqwest::get(&url).await.into_diagnostic()?;
        match response.status() {
            reqwest::StatusCode::OK => {
                let body = response.json::<LastFmResponse>().await.into_diagnostic()?;
                Ok(body.topalbums.album)
            }
            _ => Err(miette!("Error")),
        }
    }
}
