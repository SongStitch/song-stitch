# Song Stitch

<p align="center">
  <img alt="SongStitch Logo" src="https://raw.githubusercontent.com/SongStitch/song-stitch/main/public/logo.png" width="500px"/>
</p>

<div align="center">

[![Website Status](https://img.shields.io/website?style=flat-square&up_message=UP&url=https%3A%2F%2Fsongstitch.art%2F)](https://songstitch.art/)
[![CI status](https://img.shields.io/github/actions/workflow/status/SongStitch/song-stitch/deploy.yml?branch=main&style=flat-square)](https://github.com/SongStitch/song-stitch/actions?query=branch%3Amain)
[![License](https://img.shields.io/github/license/SongStitch/song-stitch?style=flat-square)](/LICENCE)

</div>

<div align="center">
A <em>blazingly fast</em> web application for generating LastFM collages, written in Go.
</div>

<br/>

<p align="center">
  <img alt="SongStitch Logo" src="https://raw.githubusercontent.com/SongStitch/song-stitch/main/docs/collage.png" width="300px"/>
</p>

## Usage

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
