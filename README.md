# Kinu

Kinu is the server that resizing images. Inspired by [cookpad's](https://cookpad.com/) image server [Tofu](http://www.slideshare.net/mirakui/ss-8150494).

## Status

It's development now for first alpha release.

## Dependency

- ImageMagick >= 6.8.9-9 && < 7.0.0

## Installation

```
# Please to Install the ImageMagick before this
go get github.com/tokubai/kinu
```

or 

```
docker pull tokubai/kinu:1.0.0.alpha4
```

## Startup

#### go get

```
$ export KINU_LOG_LEVEL=info
$ export KINU_RESIZE_ENGINE=ImageMagick
$ export KINU_STORAGE_TYPE=File 
$ export KINU_FILE_DIRECTORY=/tmp/kinu
$ kinu
INFO[0000] setup resize engine                           resize_engine_type=ImageMagick
INFO[0000] setup storage                                 storage_type=File
INFO[0000] [starting kinu on 127.0.0.1:8080]
```

```
$ curl http://localhost:8080/version
```

#### docker run

```
$ docker run takatoshimaeda/kinu:1.0.0.alpha4
INFO[0000] setup resize engine                           resize_engine_type=ImageMagick
INFO[0000] setup storage                                 storage_type=File
INFO[0000] [starting kinu on 0.0.0.0:80]
```

```
$ curl http://localhost/version
```

## Usage

### ImageUpload

#### kinu.gem(Recommended)

```ruby
require 'kinu'
Kinu::Resource.new(:foods, 1).upload(open('/path/to/image.jpg'))
```

#### curl

```shell
$ curl -X POST -F id=1 -F name=foods -F image=@/path/to/image http://localhost/upload
```

### GetResizedImage

#### kinu.gem(Recommended)

```ruby
require 'kinu'
require 'open-uri'
uri = Kinu::Resource.new(:foods, 1).uri(width: 280, height: 300) #<URI::HTTP http://localhost/images/foods/w=280,h=300/1.jpg>
open(uri).read
```

#### curl

```shell
$ curl http://localhost/images/foods/w=280,h=300/1.jpg
```

### UploadToSandbox(Temporary)

#### kinu.gem(Recommended)

```ruby
require 'kinu'
Kinu::Sandbox.upload(open('/path/to/image.jpg')) #<Kinu::Sandbox:0x007fdf92bdc490 @id="db4f1509-e2f5-40a7-9944-a6b0024f2a24", @name="__sandbox__">
```

#### curl

```shell
$ curl -X POST -F image=@/path/to/image http://localhost/sandbox
```

### Attach from Sandbox

#### kinu.gem(Recommended)

```ruby
require 'kinu'
Kinu::Sandbox.upload(open('/path/to/image.jpg')).attach_to(:foods, 1) #<Kinu::Resource:0x007fdf92b44a00 @id="1", @name="foods">
```

#### curl

```shell
$ curl -X POST -F sandbox_id=db4f1509-e2f5-40a7-9944-a6b0024f2a24 -F name=foods -F id=1 http://localhost/sandbox
```

## Specification

### Endpoints

now writing

### Image resize query

now writing

### Environment variables

| name                           | required | default value               | valid value type                                                                      | note                                                                               |
| ------------------------------ | -------- | --------------------------- | ------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| KINU_BIND                      | ☓        | 127.0.0.1:80                | IP:PORT / unix domain socket path / Einhorn(einhorn@[num]) / FileDescripter(fd@[num]) | Compliance with the specifications of the goji/bind package.                       |
| KINU_DEBUG                     | ☓        | none                        | true                                                                                  | enable pprof                                                                       |
| KINU_RESIZE_ENGINE             | ◯        | none                        | ImageMagick                                                                           |                                                                                    |
| KINU_LOG_LEVEL                 | ◯        | none                        | panic / fatal / error / warning / info / debug                                        |                                                                                    |
| KINU_LOG_FORMAT                | ☓        | text                        | ltsv / json / text                                                                    |                                                                                    |
| KINU_RESIZE_WORKER_MODE        | ☓        | none                        | true                                                                                  |                                                                                    |
| KINU_RESIZE_WORKER_MAX_SIZE    | ☓        | cpu num * 10                | Integer                                                                               |                                                                                    |
| KINU_RESIZE_WORKER_WAIT_BUFFER | ☓        | KINU_RESIZE_WORKER_SIZE * 3 | Integer                                                                               |                                                                                    |
| KINU_STORAGE_TYPE              | ◯        | none                        | File / S3                                                                             |                                                                                    |
| KINU_FILE_DIRECTORY            | ☓        | none                        | directory path                                                                        | When the `File` has been set in a `KINU_STORAGE_TYPE`\ you must set this variable. |
| KINU_S3_REGION                 | ☓        | none                        | AWS Region                                                                            | When the `S3` has been set in a `KINU_STORAGE_TYPE`\ you must set this variable.   |
| KINU_S3_BUCKET                 | ☓        | none                        | Amazon S3 bucket                                                                      | When the `S3` has been set in a `KINU_STORAGE_TYPE`\ you must set this variable.   |
| KINU_S3_BUCKET_BASE_PATH       | ☓        | none                        |                                                                                       |                                                                                    |
| AWS_ACCESS_KEY_ID              | △        | none                        |                                                                                       | Compliance with the specifications of the aws-sdk-go package.                      |
| AWS_SECRET_ACCESS_KEY          | △        | none                        |                                                                                       | Compliance with the specifications of the aws-sdk-go package.                      |

### Directory structure of the image storage.

now writing

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/tokubai/kinu

## Thanks

- [Issei Naruta](https://github.com/mirakui)
  - The [Tofu](http://www.slideshare.net/mirakui/ss-8150494), which inspired Kinu, is created by him.

## Author

- [Takatoshi Maeda](https://github.com/TakatoshiMaeda)
- [Shinya Watanabe](https://github.com/wata-gh)
