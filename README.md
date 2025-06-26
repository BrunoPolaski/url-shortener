# URL Redirection Service
This project is meant to be used by other APIs to shorten links as they need.

Please contact the project admin if you need to use it, as you will need to enter an Api-Key with its secret inside the header for link shortening.

This API is built with:

![Docker](https://img.shields.io/badge/docker-black.svg?style=for-the-badge&logo=docker&logoColor=blue)
![Go](https://img.shields.io/badge/go-black.svg?style=for-the-badge&logo=go&logoColor=blue)

Any ideas of improvement are welcomed.

Sample request to shorten a link:
```shell
curl --request POST \
  --url http://localhost:8080/link \
  --header 'Authorization: Bearer <your-bearer-token-here>' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/10.3.1' \
  --data '{
	"url": "https://my-bucket.s3.amazonaws.com/sample-object.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIOSFODNN7EXAMPLE%2F20250505%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250505T120000Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=fakeexamplesignature1234567890abcdef1234567890abcdef1234567890abcdef1234"
}'
```
>[!NOTE]
> The purpose of this shortener was to originally short S3 temporary links, this is why it is being used in the example

After this, you can just use the UUID generated along with the API route like this:

```shell
curl --request GET \
  --url http://localhost:8080/<generated-uuid>
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/10.3.1'
```
