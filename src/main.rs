mod constants;
mod image;
mod lastfm;

use ::image::ImageFormat;
use axum::{
    extract::Query,
    http::{self, StatusCode},
    response::{IntoResponse, Response},
    routing::get,
    Router,
};
use constants::{Method, Period};
use dotenv::dotenv;
use image::create_collage;
use lastfm::LastFm;
use serde::Deserialize;
use std::{io::Cursor, net::SocketAddr};
use tower_http::services::ServeDir;

#[tokio::main]
async fn main() {
    dotenv().ok();
    let app = Router::new()
        .route("/collage", get(collage))
        .nest_service("/", ServeDir::new("./public"));

    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

#[derive(Deserialize)]
struct Params {
    username: String,
    method: Method,
    period: Period,
    columns: usize,
    rows: usize,
    #[serde(rename = "album", default = "bool::default")]
    display_album: bool,
    #[serde(rename = "artist", default = "bool::default")]
    display_artist: bool,
    #[serde(rename = "track", default = "bool::default")]
    display_track: bool,
    #[serde(rename = "playcount", default = "bool::default")]
    display_playcount: bool,
    #[serde(default = "bool::default")]
    webp: bool,
}

async fn collage(Query(params): Query<Params>) -> Result<impl IntoResponse, StatusCode> {
    let lastfm = LastFm::new();
    let top_albums = lastfm
        .get_top_albums(
            &params.username,
            &params.period.to_string(),
            params.columns * params.rows,
        )
        .await
        .unwrap();
    let collage = create_collage(top_albums, params.rows, params.columns).await;
    let mut image_data = Cursor::new(Vec::new());
    let header: String;
    match params.webp {
        true => {
            collage
                .write_to(&mut image_data, ImageFormat::WebP)
                .unwrap();
            header = "image/webp".to_string();
        }
        false => {
            collage.write_to(&mut image_data, ImageFormat::Png).unwrap();
            header = "image/png".to_string()
        }
    }
    // Create the response with the image data and headers
    let response = Response::builder()
        .status(StatusCode::OK)
        .header(http::header::CONTENT_TYPE, header)
        .body(axum::body::Full::from(image_data.into_inner()))
        .unwrap();

    Ok(response)
}
