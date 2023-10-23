<a name="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
<h3 align="center">Grafana dashboards to Datadog</h3>

  <p align="center">
    This is a command line tool to convert Grafana dashboards to Datadog dashboards
    <br />
    <a href="https://github.com/abruneau/grafana_to_datadog/issues">Report Bug</a>
    ·
    <a href="https://github.com/abruneau/grafana_to_datadog/issues">Request Feature</a>
  </p>
</div>

<!-- ABOUT THE PROJECT -->
## About The Project
### Built With

* [Go](https://go.dev/)
* [goreleaser](https://goreleaser.com/)
* [cobra](https://github.com/spf13/cobra)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Brew

```sh
brew install abruneau/homebrew-tap/grafana_to_datadog
```



### From Binary

Get the [latest release](https://github.com/abruneau/grafana_to_datadog/releases) for your platform

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Build from source

1. Clone the repository
   ```sh
   git clone git@github.com:abruneau/grafana_to_datadog.git
   ```
2. Build
    ```sh
    go build -o grafana_to_datadog 
    ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

```sh
Usage:
  grafana_to_datadog [flags]

Flags:
  -d, --debug           ennable debug mode
  -h, --help            help for grafana_to_datadog
  -i, --input string    name of the input file or directory
  -o, --output string   name of the output directory
```

Convert a specific dashboard:

```sh
grafana_to_datadog -i ./my_grafana_dashboard.json -o ./output
```

Convert all dashboards in a directory 

```sh
grafana_to_datadog -i ./my_grafana_dashboard_directory -o ./output
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Compatibility

| Graphana           | Datadog    | Cloudwatch | Stackdriver | Azure |
| ------------------ | ---------- | ---------- | ----------- | ----- |
| timeseries / graph | Timeseries | ✅          | ✅           | ✅     |
| text               | Note       | ✅          | ✅           | ✅     |
| row                | Group      | ✅          | ✅           | ✅     |
| stat               | QueryValue | ✅          | ✅           | ✅     |


<!-- ROADMAP -->
## Roadmap

- [ ] Add support for more widget
- [x] Add support for GCP Stackdriver


See the [open issues](https://github.com/abruneau/grafana_to_datadog/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Your Name - antonin.bruneau@gmail.com

Project Link: [https://github.com/abruneau/grafana_to_datadog](https://github.com/abruneau/grafana_to_datadog)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/abruneau/grafana_to_datadog.svg?style=for-the-badge
[contributors-url]: https://github.com/abruneau/grafana_to_datadog/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/abruneau/grafana_to_datadog.svg?style=for-the-badge
[forks-url]: https://github.com/abruneau/grafana_to_datadog/network/members
[stars-shield]: https://img.shields.io/github/stars/abruneau/grafana_to_datadog.svg?style=for-the-badge
[stars-url]: https://github.com/abruneau/grafana_to_datadog/stargazers
[issues-shield]: https://img.shields.io/github/issues/abruneau/grafana_to_datadog.svg?style=for-the-badge
[issues-url]: https://github.com/abruneau/grafana_to_datadog/issues
[license-shield]: https://img.shields.io/github/license/abruneau/grafana_to_datadog.svg?style=for-the-badge
[license-url]: https://github.com/abruneau/grafana_to_datadog/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/antoninbruneau