# Image Compress

A simple utility program to compress image

## Supported extension

```
jpeg, svg, png, webp, heif, avif, gif
```

## Setup project

- Install `golang` in your device environment
- Clone the repo in your local environment
- Make sure to change the directory to root of the project
- Enter command in your terminal `go mod tidy`. It basically install all the required dependencies based in go.mod file
- Finally run `go run main.go` to start development server

## API Request

```
Method : POST
Endpoint : {YOUR_HOST}/api/image
Content-Type: multimedia/form-data
Body :
    image (File)
Params:
    quality (Valid integer)
    extension (jpeg | svg | png | webp | heif |avif | gif)
```
