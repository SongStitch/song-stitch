use std::fmt::Display;

use serde::Deserialize;

#[derive(Deserialize)]
pub enum Method {
    #[serde(rename = "album")]
    Album,
    #[serde(rename = "artist")]
    Artist,
    #[serde(rename = "track")]
    Track,
}

impl Display for Method {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let method = match self {
            Method::Album => "album",
            Method::Artist => "artist",
            Method::Track => "track",
        };
        write!(f, "{}", method)
    }
}

#[derive(Deserialize)]
pub enum Period {
    #[serde(rename = "7day")]
    SevenDays,
    #[serde(rename = "1month")]
    OneMonth,
    #[serde(rename = "3month")]
    ThreeMonths,
    #[serde(rename = "6month")]
    SixMonths,
    #[serde(rename = "12month")]
    TwelveMonths,
    #[serde(rename = "overall")]
    Overall,
}

impl Display for Period {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let period = match self {
            Period::SevenDays => "7day",
            Period::OneMonth => "1month",
            Period::ThreeMonths => "3month",
            Period::SixMonths => "6month",
            Period::TwelveMonths => "12month",
            Period::Overall => "overall",
        };
        write!(f, "{}", period)
    }
}
