use crate::lastfm::Album;
use ::futures::future::join_all;
use image::{imageops, DynamicImage, ImageBuffer, Rgb};
use imageproc::drawing::draw_text_mut;
use once_cell::sync::Lazy;
use rusttype::{Font, Scale};
use tokio::task;

static FONT: Lazy<Font> =
    Lazy::new(|| Font::try_from_bytes(include_bytes!("../assets/NotoSans-Regular.ttf")).unwrap());
static FONT_BOLD: Lazy<Font> =
    Lazy::new(|| Font::try_from_bytes(include_bytes!("../assets/NotoSans-Bold.ttf")).unwrap());

fn draw_text(image: &mut ImageBuffer<Rgb<u8>, Vec<u8>>, album: &Album, x: f64, y: f64) {
    let font_size = 18.0;
    let scale = Scale::uniform(font_size);
    let bg_colour = Rgb([0u8, 0u8, 0u8]);
    let fg_colour = Rgb([255u8, 255u8, 255u8]);
    draw_text_mut(
        image,
        bg_colour,
        x as i32 + 11,
        y as i32 + 11,
        scale,
        &FONT,
        &album.artist.name,
    );
    draw_text_mut(
        image,
        fg_colour,
        x as i32 + 10,
        y as i32 + 10,
        scale,
        &FONT,
        &album.artist.name,
    );
    draw_text_mut(
        image,
        bg_colour,
        x as i32 + 11,
        y as i32 + 30,
        scale,
        &FONT,
        &album.name,
    );
    draw_text_mut(
        image,
        fg_colour,
        x as i32 + 10,
        y as i32 + 29,
        scale,
        &FONT,
        &album.name,
    );
    draw_text_mut(
        image,
        bg_colour,
        x as i32 + 11,
        y as i32 + 49,
        scale,
        &FONT,
        &album.playcount,
    );
    draw_text_mut(
        image,
        fg_colour,
        x as i32 + 10,
        y as i32 + 48,
        scale,
        &FONT,
        &album.playcount,
    );
}

pub async fn create_collage(albums: Vec<Album>, rows: usize, columns: usize) -> DynamicImage {
    let mut collage =
        ImageBuffer::<Rgb<u8>, Vec<u8>>::new((300 * columns) as u32, (300 * rows) as u32);

    let tasks = albums
        .clone()
        .into_iter()
        .map(|album| task::spawn(async move { album.get_image("extralarge").await }))
        .collect::<Vec<_>>();
    let images = join_all(tasks).await;
    for (i, image_result) in images.into_iter().enumerate() {
        let image_option = image_result.unwrap();
        let x = (i as f64 % columns as f64) * 300 as f64;
        let y = ((i as f64 / columns as f64) as usize * 300) as f64;

        if let Some(mut image) = image_option {
            println!("{} {}", i, albums[i].name);
            if image.width() != 300 || image.height() != 300 {
                image = imageops::resize(&image, 300, 300, image::imageops::FilterType::Nearest);
            }
            println!("{} {} {}", i, x, y);
            imageops::overlay(&mut collage, &image, x as i64, y as i64);
            drop(image);
        }
        draw_text(&mut collage, &albums[i], x, y);
    }
    DynamicImage::ImageRgb8(collage)
}
