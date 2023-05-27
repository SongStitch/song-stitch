# Song Stitch

<p align="center">
  <img alt="SongStitch Logo" src="https://raw.githubusercontent.com/SongStitch/song-stitch/main/public/logo.png" width="500px"/>
</p>

<div align="center">

[![Website Status](https://img.shields.io/website?style=flat-square&up_message=UP&url=https%3A%2F%2Fsongstitch.art%2F)](https://songstitch.art/)
![Go Version](https://img.shields.io/github/go-mod/go-version/SongStitch/song-stitch?style=flat-square)
[![CI status](https://img.shields.io/github/actions/workflow/status/SongStitch/song-stitch/deploy.yml?branch=main&style=flat-square)](https://github.com/SongStitch/song-stitch/actions?query=branch%3Amain)
[![License](https://img.shields.io/github/license/SongStitch/song-stitch?style=flat-square)](/LICENSE)

</div>

<div align="center">
A <em>blazingly fast</em> web application for generating LastFM collages, written in Go.
</div>

<br/>

<p align="center">
  <img alt="SongStitch Collage" src="https://raw.githubusercontent.com/SongStitch/song-stitch/main/docs/collage.png" width="300px"/>
</p>

## About

SongStitch is a free, fast and highly customisable [last.fm]("https://last.fm") collage generator that allows you to create personalised collages of your most played albums. With SongStitch you can easily generate and share your collages in any size you want, displaying only the information you want, and do so amazingly quickly. Simply go to [songstitch.art](https://songstitch.art) and enter your username to start!

### Customisation Options

- **Dimensions**: specify the exact number of rows and columns you would like within your collage.
- **Text**: Choose between adding the album name, artist name and playcount to your collage; or any combo you choose.
- **Image Size**: Know the dimensions of the collage you need? SongStitch allows you to specify the desired collage image dimensions to ensure you only get what you need. You can even enable lossy compression!

Have a suggestion on how we can make SongStitch better? Feel free to create an issue on [GitHub](https://github.com/SongStitch/song-stitch/issues/new), or submit a PR!

## Usage

Below are the instructions on how you can run SongStitch yourself. You can either run it with `go` directly, or you can run it with `docker` and `docker-compose`.

1. Clone the repository

```shell
git clone git@github.com:SongStitch/song-stitch.git
```

2. Create an API key for [last.fm](https://www.last.fm/api).

3. Add environment variables to a `.env` file in the root directory. The `.env.example` includes everything that the application requires.

4. Run the application with either `docker-compose up` or `make run`. This will start the application on port `8080`.

5. Go to `localhost:8080` and enjoy!

## Contributors

- [TheDen](https://github.com/TheDen)
- [BradLewis](https://github.com/BradLewis)
