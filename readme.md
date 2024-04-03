# 62sardine

62sardine microservice is a Golang-based file upload REST API for seamless integration and efficient file management.

Created by 62teknologi.com, perfected by Community.

<details>
<summary><b>View table of contents</b></summary>

- [Running 62sardine](#running-62sardine)
  - [Prerequisites](#prerequisites)
  - [Installation manual](#installation-manual)
- [API Usage](#api-usage)
- [Contributing](#contributing)
  - [Must Preserve Characteristic](#must-preserve-characteristic)
- [About](#about-62)
</details>

## Running 62sardine

Follow the instruction below to running 62sardine on Your local machine.

### Prerequisites

Make sure to have preinstalled this prerequisites apps before You continue to installation manual. we don't include how to install these apps below most of this prerequisites is a free apps which You can find the "How to" installation tutorial anywhere in web and different machine OS have different way to install.

- MySql
- Go

### Installation manual

This installation manual will guide You to running the binary on Your ubuntu or mac terminal.

1. Clone the repository

```
git clone https://github.com/62teknologi/62sardine.git
```

2. Change directory to the cloned repository

```
cd 62sardine
```

3. tidy

```
go mod tidy
```

4. Create .env base on .env.example

```
cp .env.example .env
```

5. setup env file like this

```
HTTP_SERVER_ADDRESS="0.0.0.0:10082"
FILESYSTEM_DISK="local"
FILESYSTEM_FOLDER=sardine-services
EXPORT_FOLDER=storage
APP_URL=http://localhost:8000
```

6. Build the binary

```
go build -v -o 62sardine main.go
```

7. Run the server

```
./62sardine
```

The API server will start running on `http://localhost:10082`. You can now interact with the API using Your preferred API client or through the command line with `curl`.

## API Usage

You can find the details of API usage in [here](/docs)

## Contributing

If You'd like to contribute to the development of the 62whale REST API, please follow these steps:

1. Fork the repository
2. Create a new branch for Your feature or bugfix
3. Commit Your changes to the branch
4. Create a pull request, describing the changes You've made

We appreciate Your contributions and will review Your pull request as soon as possible.

### Must Preserve Characteristic

- Reduce repetition
- Easy to use REST API
- Easy to setup
- Easy to Customizable
- high performance
- Robust data validation and error handling
- Well documented API endpoints

## About 62

**E.nam\Du.a**

Indonesian language; spelling: A-num\Due-wa

Origin: Enam Dua means ‘six-two’ or sixty two. It is Indonesia’s international country code (+62), that was also used as a meme word for “Indonesia” by “Indonesian internet citizen” (netizen) in social media.
